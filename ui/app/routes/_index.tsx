import type { MetaFunction } from "@remix-run/node";
import { Button } from "~/components/ui/button";

export const meta: MetaFunction = () => {
  return [
    { title: "SoundScape" },
    { name: "description", content: "Music app" },
  ];
};

export default function Index() {
  return (
    <div className="flex h-screen items-center justify-center">
      <div className="flex flex-col items-center gap-1">
        <h1>Hello world</h1>
        <Button>Continue</Button>
      </div>
    </div>
  );
}
