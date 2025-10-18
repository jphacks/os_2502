import SwiftUI

struct QRScannerSheet: View {
    @Binding var scannedCode: String
    @Environment(\.dismiss) private var dismiss

    var body: some View {
        NavigationStack {
            ZStack {
                QRCodeScannerView(scannedCode: $scannedCode)

                VStack {
                    Spacer()

                    VStack(spacing: 16) {
                        Text("QRコードをスキャン")
                            .font(.title2)
                            .fontWeight(.bold)
                            .foregroundColor(.white)

                        Text("グループの招待QRコードをカメラに映してください")
                            .font(.subheadline)
                            .foregroundColor(.white.opacity(0.8))
                            .multilineTextAlignment(.center)
                    }
                    .padding()
                    .background(Color.black.opacity(0.7))
                    .cornerRadius(16)
                    .padding(.bottom, 60)
                }
            }
            .ignoresSafeArea()
            .toolbar {
                ToolbarItem(placement: .navigationBarTrailing) {
                    Button("キャンセル") {
                        dismiss()
                    }
                    .foregroundColor(.white)
                }
            }
            .toolbarBackground(.hidden, for: .navigationBar)
        }
    }
}
