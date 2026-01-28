import React from "react";
import Header from "./Header";

interface LayoutProps {
  children: React.ReactNode;
}

const Layout: React.FC<LayoutProps> = ({ children }) => {
  return (
    <div className="app-shell">
      <Header />
      <main className="app-container py-8">{children}</main>
    </div>
  );
};

export default Layout;
