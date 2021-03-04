import { useEffect, useState } from "react";

import {
  Badge,
  Card,
  CardContent,
  makeStyles,
  Typography,
  useMediaQuery,
  withStyles,
} from "@material-ui/core";

export default function SelectingGenre({
  selectedGenres,
  setSelectedGenres,
  genre,
  doneIcon,
}) {
  const classes = useStyle();
  const matches = useMediaQuery("(max-width:783px)");

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
        margin: matches ? "2vw" : "15px",
        cursor: "pointer",
      }}
    >
      <StyledBadge
        badgeContent={doneIcon()}
        invisible={!selected}
        color="primary"
      >
        <Card
          onClick={handleSelected}
          elevation={3}
          className={classes.backgroundBugOnFirefox}
        >
          <CardContent style={{ paddingBottom: "16px" }}>
            <Typography
              color="textSecondary"
              className={classes.justGenre}
              gutterBottom
            >
              Movie Genre
            </Typography>
            <Typography className={classes.genreTitle}>{genre.name}</Typography>
          </CardContent>
        </Card>
      </StyledBadge>
    </div>
  );
}

const useStyle = makeStyles((theme) => ({
  genreTitle: {
    fontSize: "1.3rem",
  },
  justGenre: {
    fontSize: "0.85rem",
  },
  backgroundBugOnFirefox: {
    background: "none",
  },
}));

const StyledBadge = withStyles((theme) => ({
  badge: {
    right: 0,
    top: 0,
    padding: "16px 6px",
    backgroundColor: "#a69c71",
  },
}))(Badge);
