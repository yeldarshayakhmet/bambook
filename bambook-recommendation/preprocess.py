import pandas as pd
from libreco.data import split_by_ratio
import random

data = pd.read_csv(
    "../dataset/goodreads_interactions.csv",
    usecols=["user_id", "book_id", "rating", "is_read"],
    engine="pyarrow")

ids = pd.read_csv("dataset/ids.csv", engine="pyarrow")
ids = ids.squeeze()

print(len(data))
data = data.loc[data["is_read"] == 1]
data.drop("is_read", axis=1, inplace=True)
data = data.loc[data["book_id"].isin(ids)]
data.rename(columns={"user_id": "user", "book_id": "item", "rating": "label"}, inplace=True)
print(len(data))

data = data.iloc[:6666666]
print(data.iloc[random.choices(range(len(data)), k=10)])  # randomly select 10 rows
data, test_data = split_by_ratio(data, test_size=0.1)
'''
chunks = np.array_split(data, 50)

for i, chunk in enumerate(chunks):
    chunk.to_csv(f"dataset/version3/interactions_train_{i+1}.csv", index=False)
'''
data.to_csv("dataset/filtered_cut/interactions.csv", index=False)
test_data.to_csv("dataset/filtered_cut/interactions_test.csv", index=False)