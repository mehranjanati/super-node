import { Room, RoomEvent, VideoPresets } from 'livekit-client';

// We assume matrix-js-sdk types are available or use any
type MatrixClient = any;
type MatrixEvent = any;

class CallManager {
    callState = $state<'IDLE' | 'RINGING' | 'CONNECTED' | 'ENDED'>('IDLE');
    callerId = $state<string | null>(null);
    currentRoomId = $state<string | null>(null);
    matrixRoomId = $state<string | null>(null);
    token = $state<string | null>(null);
    isInitialized = $state<boolean>(false);
    
    livekitRoom: Room | null = null;
    matrixClient: MatrixClient | null = null;

    reset() {
        this.callState = 'IDLE';
        this.callerId = null;
        this.currentRoomId = null;
        this.matrixRoomId = null;
        this.token = null;
        this.livekitRoom = null;
    }

    init(client: MatrixClient) {
        this.matrixClient = client;
        this.isInitialized = true;
        
        this.matrixClient.on("Room.timeline", (event: MatrixEvent, room: any, toStartOfTimeline: boolean) => {
            if (toStartOfTimeline) return;
            
            if (event.getType() === "supernode.call.incoming") {
                const content = event.getContent();
                // Only ring if we are idle
                if (this.callState === 'IDLE') {
                    console.log("Incoming call detected:", content);
                    this.callState = 'RINGING';
                    this.callerId = content.target_user || "Unknown Caller";
                    this.currentRoomId = content.livekit_room;
                    this.token = content.token;
                    this.matrixRoomId = room.roomId;
                }
            }
        });
    }

    async acceptCall() {
        const currentToken = this.token;
        const currentRoom = this.currentRoomId;
        
        if (!currentToken || !currentRoom) {
            console.error("No token or room ID to accept");
            return;
        }

        try {
            this.livekitRoom = new Room({
                adaptiveStream: true,
                dynacast: true,
                videoCaptureDefaults: {
                    resolution: VideoPresets.h720.resolution,
                },
            });

            // Handle disconnection
            this.livekitRoom.on(RoomEvent.Disconnected, () => {
                this.callState = 'IDLE';
            });

            // Connect
            // URL should be from env or config. Using localhost for dev.
            const url = import.meta.env?.VITE_LIVEKIT_URL || "ws://localhost:7880";
            await this.livekitRoom.connect(url, currentToken);
            console.log("Connected to LiveKit Room:", this.livekitRoom.name);
            
            this.callState = 'CONNECTED';
            
            // Enable tracks
            await this.livekitRoom.localParticipant.enableCameraAndMicrophone();

            // Send Accept Event to Matrix
            const mRoomId = this.matrixRoomId;
            if (this.matrixClient && mRoomId) {
                this.matrixClient.sendEvent(mRoomId, "supernode.call.accept", {
                    livekit_room: currentRoom,
                    timestamp: Date.now()
                }, "");
            }

        } catch (error) {
            console.error("Failed to connect to LiveKit", error);
            this.callState = 'IDLE'; // Reset on failure
        }
    }

    async rejectCall() {
        // Send Reject Event
        const mRoomId = this.matrixRoomId;
        const currentRoom = this.currentRoomId;
        
        if (this.matrixClient && mRoomId) {
            this.matrixClient.sendEvent(mRoomId, "supernode.call.reject", {
                livekit_room: currentRoom,
                timestamp: Date.now()
            }, "");
        }
        
        this.reset();
    }
    
    hangup() {
        if (this.livekitRoom) {
            this.livekitRoom.disconnect();
        }
        this.reset();
    }
}

export const callManager = new CallManager();
