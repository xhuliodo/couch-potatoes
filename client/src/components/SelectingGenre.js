import { useEffect, useState } from "react";

import {
  Badge,
  Card,
  CardContent,
  Typography,
  withStyles,
} from "@material-ui/core";

export default function SelectingGenre({
  selectedGenres,
  setSelectedGenres,
  genre,
  doneIcon,
}) {
  const StyledBadge = withStyles((theme) => ({
    badge: {
      right: 0,
      top: 0,
      padding: "16px 6px",
    },
  }))(Badge);

  const [selected, setSelected] = useState(false);

  const handleSelected = () => {
    setSelected(!selected);
  };

  useEffect(() => {
    selected
      ? setSelectedGenres([...selectedGenres, `${genre.genreId}`])
      : setSelectedGenres([
          ...selectedGenres.filter((sg) => sg !== `${genre.genreId}`),
        ]);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [selected]);

  return (
    <div
      style={{
        margin: "15px 15px",
        cursor: "pointer",
      }}
    >
      <StyledBadge
        badgeContent={doneIcon()}
        invisible={!selected}
        color="primary"
      >
        <Card onClick={handleSelected} elevation={6}>
          <CardContent>
            <Typography color="textSecondary" gutterBottom>
              Movie Genre
            </Typography>
            <Typography variant="h5" component="h2">
              {genre.name}
            </Typography>
          </CardContent>
        </Card>
      </StyledBadge>
    </div>
  );
}
