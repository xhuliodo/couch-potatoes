import TinderCard from "react-tinder-card";
import "./MovieCard.css";
import { Paper } from "@material-ui/core";

export default function MovieCard({ movie: m }) {
  const handleSwipe = (dir, title) => {
    // eslint-disable-next-line default-case
    switch (dir) {
      case "left":
        console.log("you just disliked: ", title);
        break;
      case "right":
        console.log("you loved the movie: ", title);
        break;
    }
  };

  return (
    <TinderCard
    
      preventSwipe={["up", "down"]}
      onSwipe={(dir) => handleSwipe(dir, m.title)}
      className="movie_swipe"
    >
      <div
        className="movie_card"
        style={{ backgroundImage: `url(${m.posterUrl})` }}
      >
        <Paper className="movie_cardText">
          <h3>
            {m.title} <i>({m.releaseYear})</i>
          </h3>
        </Paper>
      </div>
    </TinderCard>
  );
}
