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
      label: "⚙️ Configuration",
      collapsed: false,
      items: [
        "configuration/shell",
        "configuration/alias",
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
