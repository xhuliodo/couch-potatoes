import pandas as pd

movies = pd.read_csv("movies_cleaned.csv")

year = movies['title'].str[-6:]
year = year.str[:-1]
year = year.str[1:]
movies['title'] = movies['title'].str[:-7]
movies['year'] = year


movies.to_csv("movie_year.csv", index=False)