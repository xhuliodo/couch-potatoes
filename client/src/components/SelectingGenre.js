import { useState } from "react";
import { Badge, Card, CardContent, Typography } from "@material-ui/core";
import DoneOutlineIcon from "@material-ui/icons/DoneOutline";

export default function SelectingGenre() {
  const [selected, setSelected] = useState(false);

  return (
    <div>
      <Badge
        badgeContent={selected ? <DoneOutlineIcon /> : null}
        onClick={() => {
          setSelected(!selected);
        }}
        color="secondary"
      >
        <Card>
          <CardContent>
            <Typography color="textSecondary" gutterBottom>
              Word of the Day
            </Typography>
            <Typography variant="h5" component="h2">
              benevolent
            </Typography>
            <Typography color="textSecondary">adjective</Typography>
            <Typography variant="body2" component="p">
              well meaning and kindly.
              <br />
              {'"a benevolent smile"'}
            </Typography>
          </CardContent>
        </Card>
      </Badge>
    </div>
  );
}
