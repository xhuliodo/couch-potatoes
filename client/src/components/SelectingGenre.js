import { useState } from "react";
import { Badge, Card, CardContent, Typography } from "@material-ui/core";
import DoneOutlineIcon from "@material-ui/icons/DoneOutline";
import "./SelectingGenreProvider.css";

import { useQuery, useQueryClient } from "react-query";
import { request, gql } from "graphql-request";

const useGenres = () => {
  return useQuery("genres", async () => {
    const data = await request(
      "http://localhost:4001/graphql",
      gql`
        query {
          Genre {
            _id
            name
          }
        }
      `
    );
    const { Genre } = await data;
    return Genre;
  });
};

export default function SelectingGenre() {
  const [selected, setSelected] = useState(false);

  const queryClient = useQueryClient();
  const { status, data, error, isFetching } = useGenres();

  return (
    <div>
      <Typography variant="h4">Select at least 3 genres:</Typography>
      {status === "loading" ? (
        <span>Fetching data</span>
      ) : status === "error" ? (
        <span>Error: {error.message}</span>
      ) : (
        data.map((g) => (
          <Badge
            key={g._id}
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
                  {g.name}
                </Typography>
              </CardContent>
            </Card>
          </Badge>
        ))
      )}
    </div>
  );
}
