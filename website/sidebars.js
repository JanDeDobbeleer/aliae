module.exports = {
  docs: [
    "introduction",
    {
      type: "category",
      label: "📦 Installation",
      collapsed: false,
      items: [
        "installation/windows",
        "installation/macos",
        "installation/linux",
      ],
    },
    {
      type: "category",
      label: "⚙️ Setup",
      collapsed: false,
      items: [
        "setup/configuration",
        "setup/shell",
        "setup/alias",
        "setup/env",
        "setup/path",
        "setup/link",
        "setup/script",
        "setup/if",
        "setup/templates",
        "setup/include",
      ],
    },
    {
      type: "category",
      label: "🙋🏾‍♀️ Contributing",
      collapsed: true,
      items: [
        "contributing/started",
        "contributing/git",
      ],
    },
    "faq",
    "contributors",
  ],
};
