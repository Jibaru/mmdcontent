import {
	LayoutDashboard,
	MessageSquare,
	Bot,
	FileText,
	Image,
	Volume2,
	Users,
	Settings,
	CreditCard,
	History,
	Lock,
} from "lucide-react";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";

const menuItems = [
	{ icon: LayoutDashboard, label: "Main Dashboard", active: true },
	{ icon: Bot, label: "AI Assistant", active: false },
	{ icon: MessageSquare, label: "AI Chat UI", active: false },
	{ icon: FileText, label: "AI Text Generator", active: false },
	{ icon: Image, label: "AI Image Generator", active: false },
	{ icon: Volume2, label: "AI Text to Speech", active: false },
	{ icon: Users, label: "Users List", active: false },
	{ icon: Settings, label: "Profile Settings", active: false },
	{ icon: CreditCard, label: "Subscription", active: false },
	{ icon: History, label: "History", active: false },
	{ icon: Lock, label: "Authentication", active: false },
];

export function Sidebar() {
	return (
		<div className="w-64 h-screen bg-white border-r border-border flex flex-col">
			{/* Logo */}
			<div className="p-6 border-b border-border">
				<div className="flex items-center gap-3">
					<div className="w-10 h-10 bg-black rounded-lg flex items-center justify-center">
						<span className="text-white text-xl font-bold">⚡</span>
					</div>
					<span className="font-bold text-lg">Horizon AI</span>
				</div>
			</div>

			{/* Menu Items */}
			<nav className="flex-1 p-4 space-y-1">
				{menuItems.map((item) => (
					<button
						key={item.label}
						type="button"
						className={`w-full flex items-center gap-3 px-4 py-3 rounded-lg text-sm transition-colors ${
							item.active
								? "bg-black text-white"
								: "text-muted-foreground hover:bg-muted"
						}`}
					>
						<item.icon className="w-5 h-5" />
						<span>{item.label}</span>
					</button>
				))}
			</nav>

			{/* User Profile */}
			<div className="p-4 border-t border-border">
				<div className="flex items-center gap-3 p-3 rounded-lg bg-muted">
					<Avatar>
						<AvatarFallback className="bg-yellow-400 text-black">
							AP
						</AvatarFallback>
					</Avatar>
					<div className="flex-1 min-w-0">
						<div className="flex items-center gap-2">
							<span className="text-sm font-medium">PRO Member</span>
							<Badge variant="secondary" className="text-xs bg-yellow-100 text-yellow-800 hover:bg-yellow-100">
								👑
							</Badge>
						</div>
						<p className="text-xs text-muted-foreground">Unlimited plan active</p>
					</div>
				</div>
				<div className="mt-2 flex items-center justify-between px-3">
					<div className="flex items-center gap-2">
						<Avatar className="w-6 h-6">
							<AvatarFallback className="text-xs">AP</AvatarFallback>
						</Avatar>
						<span className="text-sm">Adela Parkson</span>
					</div>
					<button type="button" className="text-muted-foreground hover:text-foreground">
						<Settings className="w-4 h-4" />
					</button>
				</div>
			</div>
		</div>
	);
}
