import { useState } from "react";
import { Badge, Card, CardContent, Typography } from "@material-ui/core";
import DoneOutlineIcon from "@material-ui/icons/DoneOutline";
import { useQuery, useQueryClient } from "react-query";
import { request, gql } from "graphql-request";

import "./SelectingGenre.css";

export default function SelectingGenre() {
  const [selected, setSelected] = useState(false);
  const useGenres = () => {
    return useQuery("genres", async () => {
      const {
        genres: { data },
      } = await request(
        "http://localhost:4001/graphql",
        gql`
          query {
            Genre {
              name
            }
          }
        `
      );
      return data;
    });
  };
  const queryClient = useQueryClient();
  const { status, data, error, isFetching } = useGenres();

  return (
    <div>
      <Typography variant="h4">Select at least 3 genres:</Typography>

      <Badge
        badgeContent={selected ? <DoneOutlineIcon /> : null}
        onClick={() => {
          setSelected(!selected);
        }}
        color="secondary"
        style={{ padding: "16px 6px!important" }}
      >
        <Card>
          <CardContent>
            <Typography color="textSecondary" gutterBottom>
              Movie Genre
            </Typography>
            <Typography variant="h5" component="h2">
              benevolent
            </Typography>
          </CardContent>
        </Card>
      </Badge>
    </div>
  );
}
