[![linkedin-profile](https://img.shields.io/badge/LinkedIn-0077B5?style=for-the-badge&logo=linkedin&logoColor=white)](https://www.linkedin.com/in/xhulio-doda-745b41164/)
[![email](https://img.shields.io/badge/Gmail-D14836?style=for-the-badge&logo=gmail&logoColor=white)](mailto:xhuliodo@gmail.com)
# Couch Potatoes: Full-stack project for providing real time movie recommendations

Most people hear of new movies either from friends or from online articles reviewing them. After getting a wind of an interesting movie, you have to do the homework researching who's the director, writer or who's in the cast and then asking friends with similar movies, who might have watched it, if the movie it's worth giving a try. To simplify this process, Couch Potatoes was made to provide personalised movie recommendations.

The personalised part was made possible using a graph database (Neo4j), following two strategies:
- **Content-based** (similarity between movies)
- **User-based** (similarity between users)

The API was build using Golang, following the **Clean Architecture principles** laid by Uncle Bob.

Interfaces were built using React and the **main component is inspired by tinder's swiping cards**.

## Getting started

After cloning the repo, all you have to do is execute **docker-compose up**. 
After that you have to wait for a couple of minutes, till the seeding process finishes (the dataset used is the Small MovieLens Latest Dataset [2021]). 
(To get the full experience, you might also want to get the thumbnails which you can find the steps in the /dataset folder. It's a crawling script that downloads the thumbnails from imdb's site.)

That's all, enjoy swiping and getting some good movie recommendations.
