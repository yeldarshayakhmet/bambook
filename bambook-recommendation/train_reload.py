import pandas as pd
from libreco.algorithms import NGCF
from libreco.data import DatasetPure, DataInfo
from libreco.data import split_by_ratio


def build_dataset(index):
    data = pd.read_csv(
        f"version3/interactions_train_{index}.csv")

    data_info = DataInfo()
    if index > 1:
        data_info = DataInfo.load("model", model_name="ngcf")
        data, data_info = DatasetPure.merge_trainset(data, data_info, merge_behavior=True)
    else:
        data, data_info = DatasetPure.build_trainset(data)

    print(data_info)

    return data_info, data


def train(index):
    data_info, train_data = build_dataset(index)

    model = NGCF(
        task="ranking",
        data_info=data_info,
        n_epochs=20,
        batch_size=4096,
        loss_type="bpr",
        lr=0.005,
        device="cuda")

    if index > 1:
        model.rebuild_model(path="model", model_name="ngcf")

    model.fit(
        train_data,
        neg_sampling=True,
        verbose=2,
        metrics=["recall", "ndcg", "precision"])
    print("training complete")

    data_info.save("model_inference", model_name="ngcf")
    model.save("model_inference", model_name="ngcf", inference_only=True)

    data_info.save("model", model_name="ngcf")
    model.save("model", model_name="ngcf", manual=True, inference_only=False)

    return model, data_info


def run(index=1):
    model, data_info = None, None
    for i in range(index, 51):
        print("current index: ", i)
        model, data_info = train(i)
    return model, data_info