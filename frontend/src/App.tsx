import { useState } from "react";
import "./App.css";
import { Sidebar } from "@/components/Sidebar";
import { MainContent } from "@/components/MainContent";

export default function App() {
	const [currentView, setCurrentView] = useState("models");

	return (
		<div className="flex h-screen overflow-hidden">
			<Sidebar currentView={currentView} onViewChange={setCurrentView} />
			<MainContent view={currentView} />
		</div>
	);
}
