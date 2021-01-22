import { useEffect, useState } from "react";

import { Badge, Card, CardContent, Typography } from "@material-ui/core";
import DoneOutlineIcon from "@material-ui/icons/DoneOutline";
import "./SelectingGenre.css";

export default function SelectingGenre({
  selectedGenres,
  setSelectedGenres,
  genre,
}) {
  const [selected, setSelected] = useState(false);

  useEffect(() => {
    selected
      ? setSelectedGenres([...selectedGenres, genre.name])
      : setSelectedGenres([
          ...selectedGenres.filter((sg) => sg !== genre.name),
        ]);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [selected]);

  return (
    <Badge
      badgeContent={selected ? <DoneOutlineIcon /> : null}
      onClick={() => {
        setSelected(!selected);
      }}
      color="primary"
      style={{
        cursor: "pointer",
        padding: "2vh 6px!important",
        margin: "2.5vh",
      }}
    >
      <Card elevation={6}>
        <CardContent>
          <Typography color="textSecondary" gutterBottom>
            Movie Genre
          </Typography>
          <Typography variant="h5" component="h2">
            {genre.name}
          </Typography>
        </CardContent>
      </Card>
    </Badge>
  );
}
