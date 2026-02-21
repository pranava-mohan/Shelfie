import { View, Text, StyleSheet, Image, Pressable } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import { useAppSelector } from "@/core/store/store";
import { useAuth } from "@/core/hooks/useAuth";
import { normalize } from "@/core/normalize";
import { Ionicons } from "@expo/vector-icons";
import { useRouter } from "expo-router";

export default function AccountScreen() {
    const { name, email, avatar } = useAppSelector((state) => state.auth);
    const { logout } = useAuth();
    const router = useRouter();

    const handleLogout = async () => {
        await logout();
        router.replace("/(auth)");
    };

    return (
        <SafeAreaView style={styles.container}>
            <View style={styles.header}>
                <Text style={styles.headerTitle}>Account</Text>
            </View>

            <View style={styles.profileContainer}>
                <Image
                    source={{
                        uri: avatar || "https://ui-avatars.com/api/?name=" + (name || "User"),
                    }}
                    style={styles.avatar}
                />
                <Text style={styles.name}>{name || "User Name"}</Text>
                <Text style={styles.email}>{email || "user@example.com"}</Text>
            </View>

            <View style={styles.actionContainer}>
                <Pressable style={styles.logoutButton} onPress={handleLogout}>
                    <Ionicons name="log-out-outline" size={24} color="#FF3B30" />
                    <Text style={styles.logoutText}>Logout</Text>
                </Pressable>
            </View>
        </SafeAreaView>
    );
}

const styles = StyleSheet.create({
    container: {
        flex: 1,
        backgroundColor: "#F5F7FA",
    },
    header: {
        padding: normalize(16),
        backgroundColor: "#FFFFFF",
        borderBottomWidth: 1,
        borderBottomColor: "#E2E8F0",
        alignItems: "center",
    },
    headerTitle: {
        fontSize: normalize(18),
        fontWeight: "600",
        color: "#1A202C",
    },
    profileContainer: {
        alignItems: "center",
        paddingVertical: normalize(32),
        backgroundColor: "#FFFFFF",
        marginBottom: normalize(16),
    },
    avatar: {
        width: normalize(100),
        height: normalize(100),
        borderRadius: normalize(50),
        marginBottom: normalize(16),
        backgroundColor: "#E2E8F0",
    },
    name: {
        fontSize: normalize(20),
        fontWeight: "700",
        color: "#2D3748",
        marginBottom: normalize(4),
    },
    email: {
        fontSize: normalize(14),
        color: "#718096",
    },
    actionContainer: {
        backgroundColor: "#FFFFFF",
        paddingHorizontal: normalize(16),
    },
    logoutButton: {
        flexDirection: "row",
        alignItems: "center",
        paddingVertical: normalize(16),
        justifyContent: "center",
    },
    logoutText: {
        marginLeft: normalize(8),
        fontSize: normalize(16),
        color: "#FF3B30",
        fontWeight: "600",
    },
});
