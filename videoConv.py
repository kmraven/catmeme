import cv2
import os
import argparse

FPS = 10
VIDEO_DIR = "input_video.mp4"
FRAME_DIR = "frames"

def extract_frames(video_path, output_folder, fps):
    cap = cv2.VideoCapture(video_path)
    video_fps = cap.get(cv2.CAP_PROP_FPS)
    if not os.path.exists(output_folder):
        os.makedirs(output_folder)
    frame_count = 1
    i = 0
    while cap.isOpened():
        ret, frame = cap.read()
        if not ret:
            break
        if i % int(video_fps / fps) == 0:
            cv2.imwrite(os.path.join(output_folder, "frame_{:04d}.png".format(frame_count)), frame)
            frame_count += 1
        i += 1

    cap.release()
    cv2.destroyAllWindows()
    print(f"video converted successfully! : {video_file}")

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Extract frames from a video file at a specified FPS and save them as jpg images.")
    parser.add_argument("video_files", nargs="+", type=str, help="Path to the input video files.")
    args = parser.parse_args()

    for video_file in args.video_files:
        output_path = os.path.join(FRAME_DIR, os.path.splitext(os.path.basename(video_file))[0])
        extract_frames(video_file, output_path, FPS)
