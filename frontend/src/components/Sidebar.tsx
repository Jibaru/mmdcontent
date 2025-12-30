import {
	LayoutDashboard,
	Box,
	Layers,
	Settings,
} from "lucide-react";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";

interface SidebarProps {
	currentView: string;
	onViewChange: (view: string) => void;
}

const menuItems = [
	{ icon: LayoutDashboard, label: "Main Dashboard", view: "dashboard" },
	{ icon: Box, label: "Models", view: "models" },
	{ icon: Layers, label: "Stages", view: "stages" },
	{ icon: Settings, label: "Settings", view: "settings" },
];

export function Sidebar({ currentView, onViewChange }: SidebarProps) {
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
						onClick={() => onViewChange(item.view)}
						className={`w-full flex items-center gap-3 px-4 py-3 rounded-lg text-sm transition-colors ${
							currentView === item.view
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
