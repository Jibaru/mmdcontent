import { useState } from "react";
import {
	LayoutDashboard,
	Box,
	Layers,
	Settings,
	Sparkles,
	Zap,
} from "lucide-react";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { GenerateAll } from "../../../../wailsjs/go/handlers/Embeddings";

interface SidebarProps {
	currentView: string;
	onViewChange: (view: string) => void;
}

const menuItems = [
	{ icon: LayoutDashboard, label: "Main Dashboard", view: "dashboard" },
	{ icon: Box, label: "Models", view: "models" },
	{ icon: Layers, label: "Stages", view: "stages" },
	{ icon: Zap, label: "Motions", view: "motions" },
];

export function Sidebar({ currentView, onViewChange }: SidebarProps) {
	const [generatingEmbeddings, setGeneratingEmbeddings] = useState(false);
	const [embeddingStatus, setEmbeddingStatus] = useState<string | null>(null);

	const handleGenerateEmbeddings = async () => {
		setGeneratingEmbeddings(true);
		setEmbeddingStatus("Generating AI embeddings...");

		try {
			await GenerateAll();
			setEmbeddingStatus("✓ Embeddings generated successfully!");

			// Clear success message after 5 seconds
			setTimeout(() => {
				setEmbeddingStatus(null);
			}, 5000);
		} catch (error) {
			console.error("Error generating embeddings:", error);
			setEmbeddingStatus("✗ Error generating embeddings. Check console.");

			// Clear error message after 10 seconds
			setTimeout(() => {
				setEmbeddingStatus(null);
			}, 10000);
		} finally {
			setGeneratingEmbeddings(false);
		}
	};

	return (
		<div className="w-64 h-screen bg-white border-r border-border flex flex-col">
			{/* Logo */}
			<div className="p-6 border-b border-border">
				<div className="flex items-center gap-3">
					<div className="w-10 h-10 bg-black rounded-lg flex items-center justify-center">
						<span className="text-white text-xl font-bold">⚡</span>
					</div>
					<span className="font-bold text-lg">MMD Content</span>
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

				{/* Divider */}
				<div className="pt-4 pb-2">
					<div className="border-t border-border" />
				</div>

				{/* Generate Embeddings Button */}
				<div className="space-y-2">
					<Button
						onClick={handleGenerateEmbeddings}
						disabled={generatingEmbeddings}
						variant="outline"
						className="w-full justify-start gap-3 h-auto py-3"
					>
						{generatingEmbeddings ? (
							<div className="animate-spin rounded-full h-5 w-5 border-b-2 border-gray-900" />
						) : (
							<Sparkles className="w-5 h-5" />
						)}
						<div className="text-left flex-1">
							<div className="text-sm font-medium">
								{generatingEmbeddings ? "Generating..." : "Generate AI Search"}
							</div>
							<div className="text-xs text-muted-foreground">
								{generatingEmbeddings ? "This may take a few minutes" : "Enable semantic search"}
							</div>
						</div>
					</Button>

					{/* Status Message */}
					{embeddingStatus && (
						<div
							className={`px-3 py-2 rounded-lg text-xs ${
								embeddingStatus.startsWith("✓")
									? "bg-green-50 text-green-700 border border-green-200"
									: embeddingStatus.startsWith("✗")
									? "bg-red-50 text-red-700 border border-red-200"
									: "bg-blue-50 text-blue-700 border border-blue-200"
							}`}
						>
							{embeddingStatus}
						</div>
					)}
				</div>
			</nav>
		</div>
	);
}
