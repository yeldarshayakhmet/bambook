--this SQL script imports csv files
--each csv file line contains a single json object
--the dataset was contained in a single json file, where objects were arranged line by line, not in an object array
--the dataset was broken down to smaller chunks and converted to csv files
--it was concluded that importing PostgreSQL as csv is more convenient and faster. 
\copy book_json FROM '/Users/ellz/source/bambook-backend/dataset/partitions/part1.csv' csv;
\copy book_json FROM '/Users/ellz/source/bambook-backend/dataset/partitions/part2.csv' csv;
\copy book_json FROM '/Users/ellz/source/bambook-backend/dataset/partitions/part3.csv' csv;
\copy book_json FROM '/Users/ellz/source/bambook-backend/dataset/partitions/part4.csv' csv;
\copy book_json FROM '/Users/ellz/source/bambook-backend/dataset/partitions/part5.csv' csv;
\copy book_json FROM '/Users/ellz/source/bambook-backend/dataset/partitions/part6.csv' csv;
\copy book_json FROM '/Users/ellz/source/bambook-backend/dataset/partitions/part7.csv' csv;
\copy book_json FROM '/Users/ellz/source/bambook-backend/dataset/partitions/part8.csv' csv;
\copy book_json FROM '/Users/ellz/source/bambook-backend/dataset/partitions/part9.csv' csv;
\copy book_json FROM '/Users/ellz/source/bambook-backend/dataset/partitions/part10.csv' csv;
\copy book_json FROM '/Users/ellz/source/bambook-backend/dataset/partitions/part11.csv' csv;
\copy book_json FROM '/Users/ellz/source/bambook-backend/dataset/partitions/part12.csv' csv;
\copy book_json FROM '/Users/ellz/source/bambook-backend/dataset/partitions/part13.csv' csv;
\copy book_json FROM '/Users/ellz/source/bambook-backend/dataset/partitions/part14.csv' csv;
\copy book_json FROM '/Users/ellz/source/bambook-backend/dataset/partitions/part15.csv' csv;
\copy book_json FROM '/Users/ellz/source/bambook-backend/dataset/partitions/part16.csv' csv;
\copy book_json FROM '/Users/ellz/source/bambook-backend/dataset/partitions/part17.csv' csv;
\copy book_json FROM '/Users/ellz/source/bambook-backend/dataset/partitions/part18.csv' csv;
\copy book_json FROM '/Users/ellz/source/bambook-backend/dataset/partitions/part19.csv' csv;
\copy book_json FROM '/Users/ellz/source/bambook-backend/dataset/partitions/part20.csv' csv;
\copy book_json FROM '/Users/ellz/source/bambook-backend/dataset/partitions/part21.csv' csv;
\copy book_json FROM '/Users/ellz/source/bambook-backend/dataset/partitions/part22.csv' csv;
\copy book_json FROM '/Users/ellz/source/bambook-backend/dataset/partitions/part23.csv' csv;
\copy book_json FROM '/Users/ellz/source/bambook-backend/dataset/partitions/part24.csv' csv;

INSERT INTO books
  SELECT (jsonb_populate_record(null::books, j.book)).*
  FROM book_json j;