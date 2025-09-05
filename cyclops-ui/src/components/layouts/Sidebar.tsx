import React, { useEffect, useState } from "react";
import { Button, Menu, MenuProps } from "antd";
import {
  AppstoreAddOutlined,
  HddOutlined,
  BugFilled,
  SnippetsOutlined,
  GithubFilled,
  ThunderboltFilled,
  DiscordOutlined,
  ApiOutlined,
  RobotOutlined,
} from "@ant-design/icons";
import { useLocation } from "react-router";
import PathConstants from "../../routes/PathConstants";
import { Link } from "react-router-dom";
import "./custom.css";
import helmLogo from "../../static/img/helm_white.png";
import cyclopsLogo from "../../static/img/cyclops_logo.png";

const SideNav = () => {
  const [openKeys, setOpenKeys] = useState<string[]>([]);
  const location = useLocation(); // from react-router-dom
  const [selectedKeys, setSelectedKeys] = useState<string>("");

  useEffect(() => {
    setSelectedKeys(location.pathname.split("/")[1]);

    if (location.pathname.startsWith(PathConstants.ADDONS_MCP_SERVER)) {
      setOpenKeys(["addons"]);
    } else {
      setOpenKeys([]);
    }
  }, [location.pathname]);

  const sidebarItems: MenuProps["items"] = [
    {
      label: <Link to={PathConstants.MODULES}>Applications</Link>,
      icon: <AppstoreAddOutlined />,
      key: "modules",
    },
    {
      label: <Link to={PathConstants.TEMPLATES}>Templates</Link>,
      icon: <SnippetsOutlined />,
      key: "templates",
    },
    // {
    //   label: <Link to={PathConstants.NODES}>Nodes</Link>,
    //   icon: <HddOutlined />,
    //   key: "nodes",
    // },
    // {
    //   label: (
    //     <Link to={PathConstants.HELM_RELEASES}>
    //       Helm releases <ThunderboltFilled style={{ color: "#ffcc66" }} />
    //     </Link>
    //   ),
    //   icon: <img alt="" style={{ height: "14px" }} src={helmLogo} />,
    //   key: "helm",
    // },
    // {
    //   label: "Addons",
    //   icon: <ApiOutlined />,
    //   key: "addons",
    //   children: [
    //     {
    //       icon: <RobotOutlined />,
    //       label: <Link to={PathConstants.ADDONS_MCP_SERVER}>MCP server</Link>,
    //       key: "addons-mcp",
    //     },
    //   ],
    // },
  ];

  const tagChangelogLink = (tag: string) => {
    if (tag === "v0.0.0") {
      return "https://github.com/andersan81/cyclops/releases";
    }

    return "https://github.com/andersan81/cyclops/releases/tag/" + tag;
  };

  return (
    <div
      style={{ display: "flex", flexDirection: "column", minHeight: "100vh" }}
    >
      <Link to={PathConstants.MODULES}>
        <div
          style={{
            height: "32px",
            margin: "0.9rem 1rem 0.6rem 2rem",
            display: "inline-flex",
          }}
        >
          <img src={cyclopsLogo} alt="Cyclops" />
        </div>
      </Link>
      <Menu
        theme="dark"
        mode="inline"
        selectedKeys={[selectedKeys]}
        items={sidebarItems}
        openKeys={openKeys}
        onOpenChange={(keys) => setOpenKeys(keys)}
      />
      <Button
        style={{ background: "transparent", margin: "auto 25px 12px 25px" }}
        icon={<BugFilled />}
        className={"reportbug"}
        href={`mailto:${window.__RUNTIME_CONFIG__.REACT_APP_SUPPORT_EMAIL}?subject=Vision Deployments Bug Report`}
      >
        <b>Report a Bug</b>
      </Button>
      {/* <center
        style={{
          color: "#FFF",
          margin: "12px",
          marginTop: "0",
          fontFamily: "Arial, sans-serif",
          fontWeight: "bold",
        }}
      >
        <Link
          className={"discordlink"}
          to={"https://discord.com/invite/8ErnK3qDb3"}
        >
          <DiscordOutlined style={{ fontSize: "20px" }} /> Join Discord
        </Link>
      </center> */}
      <center
        style={{
          color: "#FFF",
          marginBottom: "25px",
          marginTop: "12px",
        }}
      >
        <span className={"taglink"}>v1.0.0</span>
      </center>
    </div>
  );
};
export default SideNav;
