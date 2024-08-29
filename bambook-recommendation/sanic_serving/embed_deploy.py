import os
from asyncio.events import AbstractEventLoop
from pathlib import Path
from typing import List

import faiss
import numpy as np
import redis.asyncio as redis
import ujson
from libreco.algorithms import YouTubeRetrieval
from libreco.data import DataInfo
from sanic import Sanic
from sanic.exceptions import SanicException
from sanic.log import logger
from sanic.request import Request
from sanic.response import HTTPResponse, json

app = Sanic("embed-serving")
data_info = DataInfo.load("/Users/ellz/source/bambook-backend/bambook-recommendation/model_inference", model_name="ngcf")
model = YouTubeRetrieval.load("/Users/ellz/source/bambook-backend/bambook-recommendation/model_inference", model_name="ngcf", data_info=data_info)


@app.get("/embed/recommend/user/<user:str>")
async def recommend(request: Request, user: str) -> HTTPResponse:
    r: redis.Redis = app.ctx.redis
    n_rec = int(request.args.get("n_rec"))
    if not await r.hexists("user2id", user):
        raise SanicException(f"Invalid user {user} doesn't exist", status_code=400)

    logger.info(f"recommend {n_rec} items for user {user}")
    user_id = await r.hget("user2id", user)
    rec_list = await recommend_on_similar_embeds(user_id, n_rec, r)
    return json(rec_list)


@app.get("/embed/recommend/cold")
async def recommend_cold(request: Request) -> HTTPResponse:
    r: redis.Redis = app.ctx.redis
    sequence = [int(item) for item in request.args.get("seq")]
    n_rec = int(request.args.get("n_rec"))
    rec_list = await recommend_dynamic(sequence, n_rec, r, True)
    return json(rec_list)


@app.get("/embed/recommend/popular")
async def recommend_popular(request: Request) -> HTTPResponse:
    n_rec = int(request.args.get("n_rec"))
    user_recommendations = model.recommend_user(user="cold", n_rec=n_rec, cold_start="popular")
    print(user_recommendations)
    return json(user_recommendations["cold"].tolist())


async def recommend_on_similar_embeds(user_id: str, n_rec: int, r: redis.Redis) -> List[int]:
    u_consumed = set(ujson.loads(await r.hget("user_consumed", user_id)))
    user_embed = ujson.loads(await r.hget("user_embed", user_id))
    # print(f"ITEM EMBED DIMENSION: {app.ctx.faiss_index.d}")
    if len(user_embed) != app.ctx.faiss_index.d:
        raise SanicException(
                "user_embed dimension != item_embed dimension, did u load the wrong faiss index?",
                status_code=500)

    user_embed = np.array(user_embed, dtype=np.float32).reshape(1, -1)
    return await compute_recommendations(n_rec, r, u_consumed, user_embed)


async def recommend_dynamic(sequence, n_rec: int, r: redis.Redis, pad=False) -> List[int]:
    user_embed = model.dyn_user_embedding(user="cold", seq=sequence)
    return await compute_recommendations(n_rec, r, sequence, np.array(user_embed, dtype=np.float32).reshape(1, -1), pad)


async def compute_recommendations(n_rec, r, u_consumed, user_embed, pad=False):
    if pad:
        user_embed = np.pad(
            user_embed,
            (0, app.ctx.faiss_index.d - user_embed.shape[1]),
            mode='constant',
            constant_values=0)
    _, item_ids = app.ctx.faiss_index.search(user_embed, n_rec + len(u_consumed))
    rec_list = []

    for i in item_ids.flatten().tolist():
        if i not in u_consumed:
            rec_list.append(int(await r.hget("id2item", i)))
            if len(rec_list) == n_rec:
                break

    return rec_list



@app.before_server_start
async def redis_faiss_setup(app: Sanic, loop: AbstractEventLoop):
    host = os.getenv("REDIS_HOST", "localhost")
    app.ctx.redis = await redis.from_url(f"redis://{host}", decode_responses=True)
    app.ctx.faiss_index = faiss.read_index(find_index_path())


@app.after_server_stop
async def redis_close(app: Sanic):
    await app.ctx.redis.close()


def find_index_path():
    # par_dir = str(Path(os.path.realpath(__file__)).parent.parent)
    par_dir = Path(__file__).absolute().parents[1]
    for dir_path, _, files in os.walk(par_dir):
        for file in files:
            if not Path(file).is_dir() and file.startswith("faiss_index"):
                return str(Path(dir_path).joinpath(file))
    raise SanicException(f"Failed to find faiss index in {par_dir}", status_code=500)


if __name__ == "__main__":
    app.run(host="0.0.0.0", port=8000, debug=False, access_log=False)