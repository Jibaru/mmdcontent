import "./App.css";
import { Sidebar } from "@/components/Sidebar";
import { Dashboard } from "@/components/Dashboard";

export default function App() {
	return (
		<div className="flex h-screen overflow-hidden">
			<Sidebar />
			<Dashboard />
		</div>
	);
}
