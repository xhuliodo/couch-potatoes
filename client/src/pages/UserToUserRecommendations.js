import { Paper, Typography } from "@material-ui/core";

export default function UserToUserRecommendations() {
  return (
    <Paper elevation={0} style={{ height: "90vh" }}>
      <Typography style={{ textAlign: "center" }} variant="h6">
        Get recommendations based on other users
      </Typography>
    </Paper>
  );
}
