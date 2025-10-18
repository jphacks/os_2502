import SwiftUI

struct FriendListView: View {
    @Environment(\.dismiss) private var dismiss
    @Environment(\.colorScheme) var colorScheme
    @State private var viewModel = FriendListViewModel()
    @State private var showAddFriend = false

    var body: some View {
        NavigationStack {
            ZStack {
                backgroundGradient
                    .ignoresSafeArea()

                ScrollView {
                    VStack(spacing: 24) {
                        if !viewModel.getPendingRequests().isEmpty {
                            pendingRequestsSection
                        }

                        friendsSection
                    }
                    .padding(.vertical, 16)
                }
            }
            .navigationTitle("フレンド")
            .navigationBarTitleDisplayMode(.large)
            .toolbar {
                ToolbarItem(placement: .topBarLeading) {
                    Button {
                        dismiss()
                    } label: {
                        Image(systemName: "xmark")
                            .foregroundColor(textPrimaryColor)
                    }
                }

                ToolbarItem(placement: .topBarTrailing) {
                    Button {
                        showAddFriend = true
                    } label: {
                        Image(systemName: "person.badge.plus")
                            .foregroundColor(textPrimaryColor)
                    }
                }
            }
            .sheet(isPresented: $showAddFriend) {
                AddFriendView(viewModel: viewModel)
            }
        }
    }

    private var pendingRequestsSection: some View {
        VStack(alignment: .leading, spacing: 16) {
            HStack(spacing: 8) {
                ZStack {
                    Circle()
                        .fill(
                            LinearGradient(
                                colors: [.orange, .red],
                                startPoint: .topLeading,
                                endPoint: .bottomTrailing
                            )
                        )
                        .frame(width: 28, height: 28)
                    Image(systemName: "bell.fill")
                        .font(.system(size: 14))
                        .foregroundColor(.white)
                }
                Text("承認待ち")
                    .font(.title3)
                    .fontWeight(.bold)
                    .foregroundColor(textPrimaryColor)
            }
            .padding(.horizontal, 24)

            VStack(spacing: 12) {
                ForEach(viewModel.getPendingRequests(), id: \.id) { friend in
                    FriendCardView(
                        friend: friend,
                        showActions: true,
                        onAccept: {
                            viewModel.acceptFriend(friend)
                        },
                        onReject: {
                            viewModel.rejectFriend(friend)
                        }
                    )
                    .padding(.horizontal, 24)
                }
            }
        }
    }

    private var friendsSection: some View {
        VStack(alignment: .leading, spacing: 16) {
            HStack(spacing: 8) {
                ZStack {
                    Circle()
                        .fill(
                            LinearGradient(
                                colors: [.blue, .cyan],
                                startPoint: .topLeading,
                                endPoint: .bottomTrailing
                            )
                        )
                        .frame(width: 28, height: 28)
                    Image(systemName: "person.2.fill")
                        .font(.system(size: 12))
                        .foregroundColor(.white)
                }
                Text("フレンド")
                    .font(.title3)
                    .fontWeight(.bold)
                    .foregroundColor(textPrimaryColor)
                Text("(\(viewModel.getFriends().count))")
                    .font(.title3)
                    .foregroundColor(.secondary)
            }
            .padding(.horizontal, 24)

            VStack(spacing: 12) {
                ForEach(viewModel.getFriends(), id: \.id) { friend in
                    FriendCardView(friend: friend)
                        .padding(.horizontal, 24)
                }
            }
        }
    }

    private var backgroundGradient: some View {
        Group {
            if colorScheme == .dark {
                LinearGradient(
                    gradient: Gradient(colors: [
                        Color(red: 0.1, green: 0.1, blue: 0.2),
                        Color(red: 0.15, green: 0.1, blue: 0.25),
                        Color(red: 0.2, green: 0.1, blue: 0.3),
                    ]),
                    startPoint: .topLeading,
                    endPoint: .bottomTrailing
                )
            } else {
                LinearGradient(
                    gradient: Gradient(colors: [
                        Color.blue.opacity(0.3),
                        Color.purple.opacity(0.3),
                        Color.pink.opacity(0.2),
                    ]),
                    startPoint: .topLeading,
                    endPoint: .bottomTrailing
                )
            }
        }
    }

    private var textPrimaryColor: Color {
        colorScheme == .dark ? .white : .primary
    }
}

struct AddFriendView: View {
    @Environment(\.dismiss) private var dismiss
    @State private var friendName = ""
    var viewModel: FriendListViewModel

    var body: some View {
        NavigationStack {
            Form {
                Section {
                    TextField("フレンド名", text: $friendName)
                }
            }
            .navigationTitle("フレンドを追加")
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .cancellationAction) {
                    Button("キャンセル") {
                        dismiss()
                    }
                }
                ToolbarItem(placement: .confirmationAction) {
                    Button("追加") {
                        viewModel.addFriend(name: friendName, iconName: "person.circle.fill")
                        dismiss()
                    }
                    .disabled(friendName.isEmpty)
                }
            }
        }
    }
}
