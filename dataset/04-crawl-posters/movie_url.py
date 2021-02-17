import csv

row_names = ['movie_id', 'movie_title', 'year', 'imdb_id']
with open('movies.txt', 'r', encoding = "ISO-8859-1") as f:
    reader = csv.DictReader(f, fieldnames=row_names, delimiter=',')
    for row in reader:
        movie_id = row['movie_id']
        movie_title = row['movie_title']
        year=row['year']
        imdb_id = row['imdb_id']
        domain = 'http://www.imdb.com/title/tt'

        while len(imdb_id)<7:
            imdb_id='0'+imdb_id

        with open('movie_url.csv', 'a', newline='') as out_csv:
                    writer = csv.writer(out_csv, delimiter=',')
                    writer.writerow([movie_id, movie_title, year, domain+imdb_id])

