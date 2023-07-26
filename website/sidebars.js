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
        "setup/alias",
        "setup/env",
        "setup/if",
        "setup/templates",
        "setup/shell",
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
