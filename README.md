Book recommendation service done as a university capstone project

A 100K subset of the [Goodreads](https://mengtingwan.github.io/data/goodreads#overview) dataset was used to train the [YouTubeRetrieval](https://librecommender.readthedocs.io/en/latest/api/algorithms/youtube_retrieval.html#libreco.algorithms.YouTubeRetrieval.dyn_user_embedding) deep learning recommendation model used from the [LibRecommender](https://github.com/massquantity/LibRecommender) library.  
The model is inferenced by deploying it to Redis and serving it through a REST API implemented on Sanic. See [libserving](https://github.com/massquantity/LibRecommender/tree/master/libserving).  
Book metadata was imported from the dataset to a PostgreSQL database. The app REST API serving the recommendations, book and user data is written using Go's Standard library.  
