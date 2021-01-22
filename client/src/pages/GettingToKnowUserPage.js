import {useStore} from "../context/genres"

export default function GettingToKnowUserPage() {
  const genres = useStore((state) => state.genres);
  console.log(genres)
  return <div>placeholder for the recommendation based on genre</div>;
}
