import { Link } from "react-router-dom";
import AvatarPopover from "./AvatarPopover";

export default function Header({
  userName,
  avatarUrl,
  onLogout,
}: {
  userName: string;
  avatarUrl?: string;
  onLogout: () => void;
}) {
  return (
    <header className="flex items-center justify-between px-6 py-4 bg-white shadow-sm">
      <Link to="/dashboard" className="flex items-center space-x-2">
        {/* Placeholder logo */}
        <img src="/placeholder-logo.svg" alt="Logo" className="h-8 w-8" />
        <span className="text-xl font-semibold text-gray-800">conquest</span>
      </Link>

      <div className="flex items-center space-x-4">
        <span className="text-gray-700">Welcome, {userName}</span>
        <AvatarPopover avatarUrl={avatarUrl} onLogout={onLogout} userName={userName} />
      </div>
    </header>
  );
}
