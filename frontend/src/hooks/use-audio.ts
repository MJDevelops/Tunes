import { Events } from "@wailsio/runtime";
import { useState, useEffect } from "react";

export function useAudio() {
  const [pos, setPos] = useState(0);
  const [playing, setPlaying] = useState(false);
  const [source, setSource] = useState("");

  const changeSource = (src: string) => {
    setSource(src);
    setPos(0);
    setPlaying(false);
  };

  useEffect(() => {}, []);

  return {
    pos,
    playing,
    changeSource,
  };
}
