import { useState, useEffect } from "react";
import { GetImageAsBase64 } from "../../wailsjs/go/main/App";
import { Card, CardContent } from "@/components/ui/card";

interface MMDContentCardProps {
	id: string;
	name: string;
	screenshots: string[];
	description: string;
	onClick?: () => void;
}

export function MMDContentCard({
	id,
	name,
	screenshots,
	description,
	onClick,
}: MMDContentCardProps) {
	const [thumbnailBase64, setThumbnailBase64] = useState<string | null>(null);
	const [loading, setLoading] = useState(true);

	useEffect(() => {
		const loadThumbnail = async () => {
			if (screenshots.length > 0) {
				try {
					const base64Image = await GetImageAsBase64(screenshots[0]);
					setThumbnailBase64(base64Image);
				} catch (error) {
					console.error(`Failed to load image for ${id}:`, error);
				} finally {
					setLoading(false);
				}
			} else {
				setLoading(false);
			}
		};

		loadThumbnail();
	}, [id, screenshots]);

	return (
		<Card
			className={`overflow-hidden hover:shadow-lg transition-shadow ${
				onClick ? "cursor-pointer" : ""
			}`}
			onClick={onClick}
		>
			<CardContent className="p-0">
				{/* Screenshot */}
				<div className="aspect-square bg-gray-100 relative">
					{thumbnailBase64 ? (
						<img
							src={thumbnailBase64}
							alt={name}
							className="w-full h-full object-cover"
						/>
					) : loading ? (
						<div className="w-full h-full flex items-center justify-center text-muted-foreground">
							<div className="flex flex-col items-center gap-2">
								<div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900" />
								<span className="text-xs">Loading...</span>
							</div>
						</div>
					) : screenshots.length > 0 ? (
						<div className="w-full h-full flex items-center justify-center text-muted-foreground">
							<div className="text-center">
								<div className="text-sm">Failed to load</div>
								<div className="text-xs">image</div>
							</div>
						</div>
					) : (
						<div className="w-full h-full flex items-center justify-center text-muted-foreground">
							No screenshot
						</div>
					)}
					<div className="absolute top-2 left-2 bg-black/70 text-white text-xs px-2 py-1 rounded">
						ID: {id}
					</div>
					{screenshots.length > 1 && (
						<div className="absolute top-2 right-2 bg-black/70 text-white text-xs px-2 py-1 rounded">
							+{screenshots.length - 1}
						</div>
					)}
				</div>

				{/* Content Info */}
				<div className="p-4">
					<h3 className="font-semibold truncate" title={name}>
						{name}
					</h3>
					<p className="text-xs text-muted-foreground mt-1 line-clamp-2">
						{description || "No description available"}
					</p>
					<div className="mt-2 flex items-center gap-2 text-xs text-muted-foreground">
						<span>{screenshots.length} screenshot{screenshots.length !== 1 ? 's' : ''}</span>
					</div>
				</div>
			</CardContent>
		</Card>
	);
}
