type PlayerComponentProps = {
  source: string;
};

export default function PlayerComponent({ source }: PlayerComponentProps) {
  return (
    <div>
      <button onClick={() => ({})}>Play</button>
    </div>
  );
}
