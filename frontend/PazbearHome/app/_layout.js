import { Drawer } from 'expo-router/drawer';

export default function Layout() {
  return ( 
    <Drawer>
      <Drawer.Screen
        name="index" // This is the name of the page and must match the url from root
        options={{
          drawerLabel: "Home",
          title: "overview",
        }}
      />
      <Drawer.Screen
        name="user/[id]" // This is the name of the page and must match the url from root
        options={{
          drawerLabel: "Minecraft",
          title: "overview",
        }}
      />
    </Drawer>
  )
}