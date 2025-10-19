import AVFoundation
import Combine
import UIKit

/// カメラ撮影を管理するサービス
class CameraService: NSObject, ObservableObject {
    @Published var isAuthorized = false
    @Published var error: CameraError?

    private let captureSession = AVCaptureSession()
    private var photoOutput = AVCapturePhotoOutput()
    private var videoPreviewLayer: AVCaptureVideoPreviewLayer?
    private var photoContinuation: CheckedContinuation<UIImage, Error>?

    enum CameraError: Error, LocalizedError {
        case notAuthorized
        case cameraUnavailable
        case captureFailed

        var errorDescription: String? {
            switch self {
            case .notAuthorized: return "カメラへのアクセスが許可されていません"
            case .cameraUnavailable: return "カメラが利用できません"
            case .captureFailed: return "撮影に失敗しました"
            }
        }
    }

    override init() {
        super.init()
    }

    // MARK: - Authorization

    /// カメラ権限をチェック
    func checkAuthorization() async -> Bool {
        switch AVCaptureDevice.authorizationStatus(for: .video) {
        case .authorized:
            await MainActor.run { isAuthorized = true }
            return true
        case .notDetermined:
            let granted = await AVCaptureDevice.requestAccess(for: .video)
            await MainActor.run { isAuthorized = granted }
            return granted
        default:
            await MainActor.run { isAuthorized = false }
            return false
        }
    }

    // MARK: - Session Setup

    /// カメラセッションをセットアップ
    func setupSession() throws {
        print("CameraService.setupSession() started")

        guard isAuthorized else {
            print("Not authorized")
            throw CameraError.notAuthorized
        }

        print("Beginning configuration...")
        captureSession.beginConfiguration()

        // カメラデバイスを取得（フロントカメラ優先）
        print("Looking for camera device...")
        guard
            let camera = AVCaptureDevice.default(
                .builtInWideAngleCamera, for: .video, position: .front)
                ?? AVCaptureDevice.default(.builtInWideAngleCamera, for: .video, position: .back)
        else {
            captureSession.commitConfiguration()
            throw CameraError.cameraUnavailable
        }
        print("Found camera: \(camera.localizedName)")

        // カメラ入力を追加
        print("Creating camera input...")
        let input = try AVCaptureDeviceInput(device: camera)
        print("Checking if can add input...")
        if captureSession.canAddInput(input) {
            print("Adding input to session...")
            captureSession.addInput(input)
            print("Input added")
        } else {
            print("Cannot add input to session")
        }

        // 写真出力を追加
        print("Checking if can add output...")
        if captureSession.canAddOutput(photoOutput) {
            print("Adding photo output...")
            captureSession.addOutput(photoOutput)

            // 高画質設定
            photoOutput.isHighResolutionCaptureEnabled = true
            if let connection = photoOutput.connection(with: .video) {
                if connection.isVideoOrientationSupported {
                    connection.videoOrientation = .portrait
                }
            }
            print("Photo output added")
        } else {
            print("Cannot add photo output")
        }

        print("Committing configuration...")
        captureSession.commitConfiguration()
        print("CameraService.setupSession() completed")
    }

    /// セッションを開始
    func startSession() {
        print("startSession() called, isRunning: \(captureSession.isRunning)")
        if !captureSession.isRunning {
            print("Starting capture session on background queue...")
            DispatchQueue.global(qos: .userInitiated).async { [weak self] in
                guard let self = self else { return }
                print("Calling captureSession.startRunning()...")
                self.captureSession.startRunning()
                print(
                    "captureSession.startRunning() returned, isRunning: \(self.captureSession.isRunning)"
                )
            }
        } else {
            print("Session already running")
        }
    }

    /// セッションを停止
    func stopSession() {
        if captureSession.isRunning {
            DispatchQueue.global(qos: .userInitiated).async { [weak self] in
                self?.captureSession.stopRunning()
            }
        }
    }

    /// プレビューレイヤーを取得
    func getPreviewLayer() -> AVCaptureVideoPreviewLayer {
        let previewLayer = AVCaptureVideoPreviewLayer(session: captureSession)
        previewLayer.videoGravity = .resizeAspectFill
        self.videoPreviewLayer = previewLayer
        return previewLayer
    }

    // MARK: - Photo Capture

    /// 写真を撮影
    func capturePhoto() async throws -> UIImage {
        guard isAuthorized else {
            throw CameraError.notAuthorized
        }

        let settings = AVCapturePhotoSettings()
        settings.flashMode = .off

        // 高画質設定
        if photoOutput.isHighResolutionCaptureEnabled {
            settings.isHighResolutionPhotoEnabled = true
        }

        return try await withCheckedThrowingContinuation { continuation in
            self.photoContinuation = continuation
            photoOutput.capturePhoto(with: settings, delegate: self)
        }
    }
}

// MARK: - AVCapturePhotoCaptureDelegate

extension CameraService: AVCapturePhotoCaptureDelegate {
    func photoOutput(
        _ output: AVCapturePhotoOutput, didFinishProcessingPhoto photo: AVCapturePhoto,
        error: Error?
    ) {
        if let error = error {
            photoContinuation?.resume(throwing: error)
            photoContinuation = nil
            return
        }

        guard let imageData = photo.fileDataRepresentation(),
            let image = UIImage(data: imageData)
        else {
            photoContinuation?.resume(throwing: CameraError.captureFailed)
            photoContinuation = nil
            return
        }

        photoContinuation?.resume(returning: image)
        photoContinuation = nil
    }
}
