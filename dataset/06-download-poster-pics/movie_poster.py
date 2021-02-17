import csv
import urllib.request

row_names = ['movieId','posterUrl']
with open('movies.csv', 'r', newline='') as in_csv:
    reader = csv.DictReader(in_csv, fieldnames=row_names, delimiter=',')
    for row in reader:
        movieId = row['movieId']
        posterUrl = row['posterUrl']
        extension = '.jpg'
        filename = 'img/' + movieId + extension
        try:
            with urllib.request.urlopen(posterUrl) as response:
                with open(filename, 'wb') as out_image:
                    out_image.write(response.read())
        except:
            with open('movies_to_delete.csv', 'a', newline='') as out_csv:
                    writer = csv.writer(out_csv, delimiter=',')
                    writer.writerow([movieId])