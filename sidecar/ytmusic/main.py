"""YouTube Music Sidecar Service using ytmusicapi."""

import logging
from contextlib import asynccontextmanager
from typing import Optional

from fastapi import FastAPI, HTTPException, Query
from pydantic import BaseModel
from ytmusicapi import YTMusic

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s - %(name)s - %(levelname)s - %(message)s",
)
logger = logging.getLogger(__name__)

# Global ytmusic client
ytmusic: Optional[YTMusic] = None


class Track(BaseModel):
    """Track model for API responses."""

    video_id: str
    title: str
    artist: str
    artist_id: Optional[str] = None
    album: Optional[str] = None
    album_id: Optional[str] = None
    duration_seconds: Optional[int] = None
    thumbnail_url: Optional[str] = None
    is_explicit: bool = False


class SimilarTracksResponse(BaseModel):
    """Response model for similar tracks endpoint."""

    video_id: str
    tracks: list[Track]


class SearchResponse(BaseModel):
    """Response model for search endpoint."""

    query: str
    tracks: list[Track]


class HealthResponse(BaseModel):
    """Response model for health check."""

    status: str
    service: str


@asynccontextmanager
async def lifespan(app: FastAPI):
    """Initialize ytmusicapi client on startup."""
    global ytmusic
    logger.info("Initializing YTMusic client...")
    # No auth required for public data
    ytmusic = YTMusic()
    logger.info("YTMusic client initialized successfully")
    yield
    logger.info("Shutting down YTMusic sidecar...")


app = FastAPI(
    title="YouTube Music Sidecar",
    description="Sidecar service for YouTube Music API using ytmusicapi",
    version="1.0.0",
    lifespan=lifespan,
)


@app.get("/health", response_model=HealthResponse)
async def health_check():
    """Health check endpoint."""
    return HealthResponse(status="healthy", service="ytmusic-sidecar")


@app.get("/similar/{video_id}", response_model=SimilarTracksResponse)
async def get_similar_tracks(
    video_id: str,
    limit: int = Query(default=25, ge=1, le=100, description="Max number of tracks to return"),
):
    """
    Get similar/radio tracks for a given video ID.

    Uses ytmusicapi's get_watch_playlist which returns the "radio" playlist
    for a given song - these are YouTube Music's recommendations.
    """
    if ytmusic is None:
        raise HTTPException(status_code=503, detail="YTMusic client not initialized")

    try:
        logger.info(f"Getting similar tracks for video_id={video_id}, limit={limit}")

        # get_watch_playlist returns the autoplay/radio queue for a video
        watch_playlist = ytmusic.get_watch_playlist(videoId=video_id, limit=limit + 10)

        if not watch_playlist or "tracks" not in watch_playlist:
            logger.warning(f"No watch playlist found for video_id={video_id}")
            return SimilarTracksResponse(video_id=video_id, tracks=[])

        tracks = []
        for track_data in watch_playlist.get("tracks", [])[1:]:  # Skip first (seed track)
            if len(tracks) >= limit:
                break

            track = _parse_track(track_data)
            if track:
                tracks.append(track)

        logger.info(f"Found {len(tracks)} similar tracks for video_id={video_id}")
        return SimilarTracksResponse(video_id=video_id, tracks=tracks)

    except Exception as e:
        logger.error(f"Error getting similar tracks for {video_id}: {e}")
        raise HTTPException(status_code=500, detail=f"Failed to get similar tracks: {str(e)}")


@app.get("/search", response_model=SearchResponse)
async def search_tracks(
    q: str = Query(..., min_length=1, description="Search query (artist - track or track name)"),
    limit: int = Query(default=5, ge=1, le=50, description="Max number of results"),
):
    """
    Search for tracks on YouTube Music.

    Useful for finding video IDs when you only have artist/track name.
    """
    if ytmusic is None:
        raise HTTPException(status_code=503, detail="YTMusic client not initialized")

    try:
        logger.info(f"Searching for: {q}, limit={limit}")

        search_results = ytmusic.search(query=q, filter="songs", limit=limit)

        if not search_results:
            return SearchResponse(query=q, tracks=[])

        tracks = []
        for result in search_results:
            if len(tracks) >= limit:
                break

            track = _parse_search_result(result)
            if track:
                tracks.append(track)

        logger.info(f"Found {len(tracks)} tracks for query: {q}")
        return SearchResponse(query=q, tracks=tracks)

    except Exception as e:
        logger.error(f"Error searching for '{q}': {e}")
        raise HTTPException(status_code=500, detail=f"Search failed: {str(e)}")


def _parse_track(track_data: dict) -> Optional[Track]:
    """Parse track data from watch playlist response."""
    try:
        video_id = track_data.get("videoId")
        title = track_data.get("title")

        if not video_id or not title:
            return None

        # Get artist info
        artists = track_data.get("artists", [])
        artist = artists[0].get("name", "Unknown") if artists else "Unknown"
        artist_id = artists[0].get("id") if artists else None

        # Get album info
        album_data = track_data.get("album")
        album = album_data.get("name") if album_data else None
        album_id = album_data.get("id") if album_data else None

        # Duration (comes as "3:45" string, convert to seconds)
        duration_str = track_data.get("length")
        duration_seconds = _parse_duration(duration_str) if duration_str else None

        # Thumbnail
        thumbnails = track_data.get("thumbnail", [])
        thumbnail_url = thumbnails[-1].get("url") if thumbnails else None

        return Track(
            video_id=video_id,
            title=title,
            artist=artist,
            artist_id=artist_id,
            album=album,
            album_id=album_id,
            duration_seconds=duration_seconds,
            thumbnail_url=thumbnail_url,
            is_explicit=track_data.get("isExplicit", False),
        )
    except Exception as e:
        logger.warning(f"Failed to parse track: {e}")
        return None


def _parse_search_result(result: dict) -> Optional[Track]:
    """Parse track data from search result."""
    try:
        video_id = result.get("videoId")
        title = result.get("title")

        if not video_id or not title:
            return None

        # Get artist info
        artists = result.get("artists", [])
        artist = artists[0].get("name", "Unknown") if artists else "Unknown"
        artist_id = artists[0].get("id") if artists else None

        # Get album info
        album_data = result.get("album")
        album = album_data.get("name") if album_data else None
        album_id = album_data.get("id") if album_data else None

        # Duration
        duration_str = result.get("duration")
        duration_seconds = _parse_duration(duration_str) if duration_str else None

        # Thumbnail
        thumbnails = result.get("thumbnails", [])
        thumbnail_url = thumbnails[-1].get("url") if thumbnails else None

        return Track(
            video_id=video_id,
            title=title,
            artist=artist,
            artist_id=artist_id,
            album=album,
            album_id=album_id,
            duration_seconds=duration_seconds,
            thumbnail_url=thumbnail_url,
            is_explicit=result.get("isExplicit", False),
        )
    except Exception as e:
        logger.warning(f"Failed to parse search result: {e}")
        return None


def _parse_duration(duration_str: str) -> Optional[int]:
    """Parse duration string like '3:45' or '1:02:30' to seconds."""
    try:
        parts = duration_str.split(":")
        if len(parts) == 2:
            return int(parts[0]) * 60 + int(parts[1])
        elif len(parts) == 3:
            return int(parts[0]) * 3600 + int(parts[1]) * 60 + int(parts[2])
        return None
    except (ValueError, AttributeError):
        return None


if __name__ == "__main__":
    import uvicorn

    uvicorn.run(app, host="0.0.0.0", port=8081)
