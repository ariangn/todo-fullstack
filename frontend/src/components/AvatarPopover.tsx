import { Popover, PopoverTrigger, PopoverContent } from "@/components/ui/popover";
import { Button } from "@/components/ui/button";
import { User2Icon, LogOutIcon } from "lucide-react";
import { randomColorFromString } from "../utils/colorHash";

export default function AvatarPopover({
  avatarUrl,
  userName,
  onLogout,
}: {
  avatarUrl?: string;
  userName: string;
  onLogout: () => void;
}) {
  const bgColor = avatarUrl ? "bg-transparent" : randomColorFromString(userName); // e.g. "bg-red-500"

  return (
    <Popover>
      <PopoverTrigger asChild>
        <div
          className={`h-10 w-10 rounded-full overflow-hidden border border-gray-300 flex items-center justify-center ${bgColor}`}
        >
          {avatarUrl ? (
            <img
              src={avatarUrl}
              alt="Avatar"
              className="h-full w-full object-cover"
            />
          ) : (
            <User2Icon className="h-6 w-6 text-white" />
          )}
        </div>
      </PopoverTrigger>
      <PopoverContent align="end" className="w-32 p-2 bg-white">
        <Button
          variant="ghost"
          size="sm"
          className="w-full justify-start space-x-2"
          onClick={onLogout}
        >
          <LogOutIcon className="h-4 w-4" />
          <span>Sign Out</span>
        </Button>
      </PopoverContent>
    </Popover>
  );
}
