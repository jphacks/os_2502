import SwiftUI

struct JoinGroupView: View {
    let authManager: AuthenticationManager
    @State private var viewModel: CollageGroupViewModel?
    @State private var inviteCode: String = ""
    @State private var showingError = false
    @State private var errorMessage = ""
    @State private var showingWaitingRoom = false
    @State private var showingQRScanner = false
    @State private var isJoining = false
    @Environment(\.dismiss) private var dismiss

    var body: some View {
        VStack(spacing: 30) {
            Text("グループに参加")
                .font(.largeTitle)
                .fontWeight(.bold)

            VStack(alignment: .leading, spacing: 12) {
                Text("招待コード")
                    .font(.headline)

                TextField("招待コードを入力", text: $inviteCode)
                    .textFieldStyle(.roundedBorder)
                    .textInputAutocapitalization(.characters)
                    .font(.title3)
                    .padding(.horizontal)

                Text("グループ作成者から共有された招待コードを入力してください")
                    .font(.caption)
                    .foregroundColor(.gray)
                    .padding(.horizontal)
            }
            .padding()
            .background(Color.gray.opacity(0.1))
            .cornerRadius(12)

            Text("または")
                .font(.subheadline)
                .foregroundColor(.gray)

            Button {
                showingQRScanner = true
            } label: {
                HStack {
                    Image(systemName: "qrcode.viewfinder")
                        .font(.title2)
                    Text("QRコードをスキャン")
                        .font(.title3)
                        .fontWeight(.semibold)
                }
                .foregroundColor(.blue)
                .frame(maxWidth: .infinity)
                .padding()
                .background(Color.blue.opacity(0.1))
                .cornerRadius(12)
            }

            Spacer()

            Button {
                Task {
                    await joinGroup()
                }
            } label: {
                HStack {
                    if isJoining {
                        ProgressView()
                            .tint(.white)
                    }
                    Text(isJoining ? "参加中..." : "参加する")
                        .font(.title2)
                        .fontWeight(.semibold)
                }
                .foregroundColor(.white)
                .frame(maxWidth: .infinity)
                .padding()
                .background(inviteCode.isEmpty || isJoining ? Color.gray : Color.blue)
                .cornerRadius(12)
            }
            .disabled(inviteCode.isEmpty || isJoining)
        }
        .padding()
        .alert("エラー", isPresented: $showingError) {
            Button("OK", role: .cancel) {}
        } message: {
            Text(errorMessage)
        }
        .sheet(isPresented: $showingQRScanner) {
            QRScannerSheet(scannedCode: $inviteCode)
        }
        .navigationDestination(isPresented: $showingWaitingRoom) {
            if let viewModel = viewModel {
                SimpleWaitingRoomView(viewModel: viewModel)
            }
        }
        .onChange(of: inviteCode) { _, newValue in
            if !newValue.isEmpty && showingQRScanner {
                showingQRScanner = false
            }
        }
        .navigationBarBackButtonHidden(true)
        .toolbar {
            ToolbarItem(placement: .navigationBarLeading) {
                Button {
                    dismiss()
                } label: {
                    HStack {
                        Image(systemName: "xmark")
                        Text("閉じる")
                    }
                }
            }
        }
        .onAppear {
            if viewModel == nil {
                viewModel = CollageGroupViewModel(authManager: authManager)
            }
        }
    }

    private func joinGroup() async {
        guard let vm = viewModel else { return }
        guard !inviteCode.isEmpty else { return }

        isJoining = true

        // APIでグループに参加
        let success = await vm.joinGroupWithAPI(invitationToken: inviteCode)

        isJoining = false

        if success {
            showingWaitingRoom = true
        } else {
            errorMessage = vm.errorMessage ?? "グループへの参加に失敗しました"
            showingError = true
        }
    }
}

#Preview {
    NavigationStack {
        JoinGroupView(authManager: AuthenticationManager())
    }
}
