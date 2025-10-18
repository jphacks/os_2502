import SwiftUI

struct ResultView: View {
    @Bindable var viewModel: CollageGroupViewModel
    var capturedImage: UIImage?
    @Environment(\.dismiss) private var dismiss

    var body: some View {
        VStack(spacing: 30) {
            Text("コラージュ完成")
                .font(.largeTitle)
                .fontWeight(.bold)

            if let group = viewModel.currentGroup {
                Text("\(group.members.count)人のコラージュ")
                    .font(.headline)
                    .foregroundColor(.gray)
            }

            ScrollView {
                VStack(spacing: 16) {
                    if let group = viewModel.currentGroup {
                        ForEach(group.members) { member in
                            VStack(alignment: .leading, spacing: 8) {
                                HStack {
                                    Image(systemName: "person.circle.fill")
                                        .foregroundColor(.blue)
                                    Text(member.name)
                                        .font(.headline)
                                    Spacer()
                                }

                                // 撮影された画像のプレースホルダー
                                RoundedRectangle(cornerRadius: 12)
                                    .fill(Color.gray.opacity(0.3))
                                    .frame(height: 200)
                                    .overlay(
                                        Image(systemName: "photo")
                                            .font(.system(size: 60))
                                            .foregroundColor(.gray)
                                    )
                            }
                            .padding()
                            .background(Color.white)
                            .cornerRadius(12)
                            .shadow(radius: 2)
                        }
                    }
                }
                .padding(.horizontal)
            }

            Spacer()

            VStack(spacing: 16) {
                Button {
                    // コラージュを保存
                } label: {
                    HStack {
                        Image(systemName: "square.and.arrow.down")
                        Text("コラージュを保存")
                    }
                    .font(.title3)
                    .fontWeight(.semibold)
                    .foregroundColor(.white)
                    .frame(maxWidth: .infinity)
                    .padding()
                    .background(Color.blue)
                    .cornerRadius(12)
                }

                Button {
                    viewModel.resetGroup()
                    dismiss()
                } label: {
                    Text("終了")
                        .font(.title3)
                        .fontWeight(.semibold)
                        .foregroundColor(.blue)
                        .frame(maxWidth: .infinity)
                        .padding()
                        .background(Color.gray.opacity(0.2))
                        .cornerRadius(12)
                }
            }
            .padding(.horizontal)
        }
        .padding()
        .navigationBarBackButtonHidden(true)
    }
}

#Preview {
    let authManager = AuthenticationManager()
    let viewModel = CollageGroupViewModel(authManager: authManager)
    let _ = viewModel.createGroupLocal(type: .temporaryLocal, maxMembers: 3)

    return NavigationStack {
        ResultView(viewModel: viewModel, capturedImage: nil)
    }
}
