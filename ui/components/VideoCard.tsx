import Image from "next/image";
import Link from "next/link";
import { Calendar, ExternalLink } from "lucide-react";
import { useState } from "react";

interface Video {
    vid: string
    title: string
    description: string
    published_at: string
    thumbnail: string
}

export default function VideoCard({ video }: { video: Video }) {
    const [imageError, setImageError] = useState(false);

    return (
        <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md overflow-hidden transition-transform duration-200 hover:scale-105">
            <Link
                href={`https://www.youtube.com/watch?v=${video.vid}`}
                target="_blank"
                rel="noopener noreferrer"
                className="block"
            >
                <div className="relative">
                    <Image
                        src={imageError ? "/placeholder.svg" : video.thumbnail}
                        alt={video.title}
                        width={200}
                        height={112}
                        className="w-full h-48 object-cover"
                        onError={() => setImageError(true)}
                    />
                    <div className="absolute top-2 right-2 bg-gray-900 bg-opacity-75 text-white p-1 rounded">
                        <ExternalLink size={16} />
                    </div>
                </div>
                <div className="p-4">
                    <h2 className="text-lg font-semibold text-gray-800 dark:text-gray-100 mb-2 line-clamp-2">
                        {video.title}
                    </h2>
                    <p className="text-sm text-gray-600 dark:text-gray-300 mb-2 line-clamp-3">
                        {video.description}
                    </p>
                    <div className="flex items-center text-xs text-gray-500 dark:text-gray-400">
                        <Calendar size={14} className="mr-1" />
                        {new Date(video.published_at).toLocaleDateString()}
                    </div>
                </div>
            </Link>
        </div>
    );
}
