import pandas as pd
from libreco.data import DatasetPure
from libreco.algorithms import YouTubeRetrieval
import sys

limit = sys.getrecursionlimit()
print(limit)
sys.setrecursionlimit(limit * 2)
print(sys.getrecursionlimit())

data = pd.read_csv(f"dataset/filtered_cut/interactions.csv")
data, data_info = DatasetPure.build_trainset(data)
print(data_info)

model = YouTubeRetrieval(
    task="ranking",
    data_info=data_info,
    n_epochs=30,
    batch_size=1024)

model.fit(data, neg_sampling=True, shuffle=False, verbose=2)

data_info.save("model_inference", model_name="ngcf")
model.save("model_inference", model_name="ngcf", inference_only=True)

data_info.save("model", model_name="ngcf")
model.save("model", model_name="ngcf", manual=True, inference_only=False)