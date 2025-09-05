import { ExclamationCircleTwoTone } from "@ant-design/icons";
import React, { CSSProperties } from "react";
import { useTheme } from "../../theme/ThemeContext";

interface ErrorIconProps {
  style?: CSSProperties;
}

export function ErrorIcon({ style }: ErrorIconProps) {
  const { mode } = useTheme();

  return (
    <ExclamationCircleTwoTone
      style={style}
      twoToneColor={mode === "light" ? "#ffcc00" : ["#ffcc00", "#4f4002"]}
    />
  );
}
