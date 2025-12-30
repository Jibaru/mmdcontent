import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
	Table,
	TableBody,
	TableCell,
	TableHead,
	TableHeader,
	TableRow,
} from "@/components/ui/table";
import { BarChart3, Users, CreditCard, DollarSign } from "lucide-react";

const statsCards = [
	{
		title: "Total Credits Used",
		value: "46,823",
		change: "+20.4% from last month",
		icon: BarChart3,
	},
	{
		title: "Total Users",
		value: "67,284",
		change: "+12.3% from last month",
		icon: Users,
	},
	{
		title: "Credits Available",
		value: "100,000",
		change: "",
		icon: CreditCard,
	},
	{
		title: "Current Expenses",
		value: "Expected",
		change: "",
		icon: DollarSign,
	},
];

const userData = [
	{
		email: "hello@horizon-ui.com",
		provider: "Google",
		created: "06 Nov, 2023 11:33",
		lastSignIn: "06 Nov, 2023 11:33",
		userId: "f3f42fc419-ce32-49fc-92df...",
	},
	{
		email: "thomas@gmail.com",
		provider: "Google",
		created: "06 Nov, 2023 11:29",
		lastSignIn: "06 Nov, 2023 11:29",
		userId: "f3f42fc419-ce32-49fc-92df...",
	},
	{
		email: "markwilliam@hotmail.com",
		provider: "Email",
		created: "06 Nov, 2023 11:21",
		lastSignIn: "06 Nov, 2023 11:21",
		userId: "f3f42fc419-ce32-49fc-92df...",
	},
	{
		email: "examplejosh@mail.com",
		provider: "Google",
		created: "06 Nov, 2023 11:18",
		lastSignIn: "06 Nov, 2023 11:18",
		userId: "f3f42fc419-ce32-49fc-92df...",
	},
];

const months = [
	"SEP",
	"OCT",
	"NOV",
	"DEC",
	"JAN",
	"FEB",
	"MAR",
	"APR",
	"MAY",
	"JUN",
];

export function Dashboard() {
	return (
		<div className="flex-1 overflow-auto bg-gray-50">
			{/* Header */}
			<div className="bg-white border-b border-border px-8 py-4">
				<div className="flex items-center gap-2 text-sm text-muted-foreground mb-2">
					<span>Pages</span>
					<span>/</span>
					<span>Main Dashboard</span>
				</div>
				<h1 className="text-3xl font-bold">Main Dashboard</h1>
			</div>

			{/* Content */}
			<div className="p-8 space-y-6">
				{/* Stats Cards */}
				<div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
					{statsCards.map((stat) => (
						<Card key={stat.title}>
							<CardContent className="p-6">
								<div className="flex items-start justify-between">
									<div>
										<p className="text-sm text-muted-foreground mb-1">
											{stat.title}
										</p>
										<h3 className="text-3xl font-bold">{stat.value}</h3>
										{stat.change && (
											<p className="text-xs text-green-600 mt-2">
												{stat.change}
											</p>
										)}
									</div>
									<div className="w-12 h-12 bg-gray-100 rounded-lg flex items-center justify-center">
										<stat.icon className="w-6 h-6 text-gray-600" />
									</div>
								</div>
							</CardContent>
						</Card>
					))}
				</div>

				{/* Credits Usage Chart */}
				<Card>
					<CardHeader>
						<div className="flex items-start justify-between">
							<div>
								<CardTitle>Credits Usage Last Year</CardTitle>
								<p className="text-3xl font-bold mt-2">249,758</p>
							</div>
							<div className="w-12 h-12 bg-gray-100 rounded-lg flex items-center justify-center">
								<BarChart3 className="w-6 h-6 text-gray-600" />
							</div>
						</div>
					</CardHeader>
					<CardContent>
						{/* Simple wave chart placeholder */}
						<div className="relative h-64">
							<svg
								className="w-full h-full"
								viewBox="0 0 1200 250"
								preserveAspectRatio="none"
							>
								<path
									d="M 0,125 Q 150,50 300,100 T 600,125 T 900,75 T 1200,50"
									fill="none"
									stroke="currentColor"
									strokeWidth="2"
									className="text-black"
								/>
							</svg>
							{/* Month labels */}
							<div className="flex justify-between mt-4 text-xs text-muted-foreground">
								{months.map((month) => (
									<span key={month}>{month}</span>
								))}
							</div>
						</div>
					</CardContent>
				</Card>

				{/* Users Table */}
				<Card>
					<CardContent className="p-6">
						<Table>
							<TableHeader>
								<TableRow>
									<TableHead className="w-12">
										<input type="checkbox" className="rounded border-gray-300" />
									</TableHead>
									<TableHead>EMAIL ADDRESS</TableHead>
									<TableHead>PROVIDER</TableHead>
									<TableHead>CREATED</TableHead>
									<TableHead>LAST SIGN IN</TableHead>
									<TableHead>USER UID</TableHead>
								</TableRow>
							</TableHeader>
							<TableBody>
								{userData.map((user) => (
									<TableRow key={user.email}>
										<TableCell>
											<input
												type="checkbox"
												className="rounded border-gray-300"
											/>
										</TableCell>
										<TableCell className="font-medium">{user.email}</TableCell>
										<TableCell>{user.provider}</TableCell>
										<TableCell>{user.created}</TableCell>
										<TableCell>{user.lastSignIn}</TableCell>
										<TableCell className="text-muted-foreground">
											{user.userId}
										</TableCell>
									</TableRow>
								))}
							</TableBody>
						</Table>
						<div className="mt-4 text-sm text-muted-foreground">
							0 of 4 row(s) selected.
						</div>
					</CardContent>
				</Card>
			</div>
		</div>
	);
}
