import { ModelsGrid } from "../../screens/ModelsGrid";

interface MainContentProps {
	view: string;
}

export function MainContent({ view }: MainContentProps) {
	return (
		<div className="flex-1 overflow-auto bg-gray-50">
			{/* Header */}
			<div className="bg-white border-b border-border px-8 py-4">
				<div className="flex items-center gap-2 text-sm text-muted-foreground mb-2">
					<span>Pages</span>
					<span>/</span>
					<span>{view === "models" ? "Models" : view === "stages" ? "Stages" : "Dashboard"}</span>
				</div>
				<h1 className="text-3xl font-bold">
					{view === "models" ? "Models Library" : view === "stages" ? "Stages Library" : "Main Dashboard"}
				</h1>
			</div>

			{/* Content */}
			<div className="p-8">
				{view === "models" && <ModelsGrid />}
				{view === "stages" && (
					<div className="text-center text-muted-foreground py-12">
						Stages view coming soon...
					</div>
				)}
				{view === "dashboard" && (
					<div className="text-center text-muted-foreground py-12">
						Dashboard statistics coming soon...
					</div>
				)}
			</div>
		</div>
	);
}
