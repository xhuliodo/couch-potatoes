import { useEffect, useState } from "react";

// Tracks a poster image that failed to load, resetting whenever the URL
// changes so a later card with a working image shows its poster.
export default function usePosterFallback(image) {
  const [failed, setFailed] = useState(!image);

  useEffect(() => {
    setFailed(!image);
  }, [image]);

  return {
    showFallback: failed,
    onError: () => setFailed(true),
  };
}
