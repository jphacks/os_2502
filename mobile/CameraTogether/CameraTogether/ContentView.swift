import SwiftData
import SwiftUI

struct ContentView: View {
    @Environment(\.modelContext) private var modelContext
    @Environment(\.colorScheme) var colorScheme
    @State private var isShowingSettings = false
    @State private var authManager = AuthenticationManager()

    var body: some View {
        if !authManager.isAuthenticated {
            LoginView(authManager: authManager)
        } else {
            mainContent
        }
    }

    private var mainContent: some View {
        NavigationStack {
            ZStack {
                // TODO: ここに背景
                ScrollView {
                }
            }
            .navigationTitle("Collage")
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .navigationBarTrailing) {
                    Button {
                        isShowingSettings = true
                    } label: {
                        Image(systemName: "gearshape.fill")
                            .foregroundColor(textPrimaryColor)
                    }
                }
            }
            .toolbarBackground(.visible, for: .navigationBar)
            .toolbarBackground(Color.clear, for: .navigationBar)
            .sheet(isPresented: $isShowingSettings) {
                SettingsSheetView(authManager: authManager)
            }
        }
    }
    private var textPrimaryColor: Color {
        colorScheme == .dark ? .white : .primary
    }
}

#Preview {
    ContentView()
        .modelContainer(for: Item.self, inMemory: true)
}
