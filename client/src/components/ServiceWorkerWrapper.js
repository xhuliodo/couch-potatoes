import React, { useEffect, useState } from "react";
import { Snackbar, Button } from "@material-ui/core";
// pwa service worker
import * as serviceWorker from "../serviceWorkerRegistration";

const ServiceWorkerWrapper = () => {
  const [showReload, setShowReload] = useState(false);
  const [waitingWorker, setWaitingWorker] = useState(null);

  const onSWUpdate = (registration) => {
    setShowReload(true);
    setWaitingWorker(registration.waiting);
  };

  useEffect(() => {
    serviceWorker.register({ onUpdate: onSWUpdate });
  }, []);

  const reloadPage = () => {
    waitingWorker?.postMessage({ type: "SKIP_WAITING" });
    setShowReload(false);
    window.location.reload();
  };

  return (
    <Snackbar
      open={showReload}
      message="A new version is available!"
      onClick={reloadPage}
      anchorOrigin={{ vertical: "bottom", horizontal: "center" }}
      action={
        <Button
          variant="contained"
          size="small"
          onClick={reloadPage}
        >
          Update
        </Button>
      }
    />
  );
};

export default ServiceWorkerWrapper;
