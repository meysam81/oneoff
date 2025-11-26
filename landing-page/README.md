# OneOff Landing Page

The marketing and documentation website for [OneOff](https://github.com/meysam81/oneoff) — a self-hosted one-time job scheduler.

## Tech Stack

- **Framework**: [Astro](https://astro.build/) v5
- **UI Components**: [Vue 3](https://vuejs.org/) (for interactive islands)
- **Styling**: [Tailwind CSS](https://tailwindcss.com/)
- **Deployment**: GitHub Pages

## Development

### Prerequisites

- Node.js 20+
- npm or pnpm

### Setup

```bash
# Install dependencies
npm install

# Start development server
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview
```

The development server runs at `http://localhost:4321`.

### Project Structure

```
landing-page/
├── src/
│   ├── components/
│   │   ├── ui/          # Reusable UI primitives
│   │   ├── landing/     # Landing page sections
│   │   ├── catalog/     # Catalog Vue components
│   │   └── icons/       # SVG icon component
│   ├── layouts/         # Page layouts
│   ├── pages/           # Route pages
│   ├── content/
│   │   └── catalog/     # Job template JSON files
│   └── styles/          # Global CSS
├── public/              # Static assets
└── astro.config.mjs     # Astro configuration
```

## Contributing Templates

To add a job template to the catalog:

1. Create a JSON file in `src/content/catalog/`
2. Follow the template schema:

```json
{
  "id": "unique-template-id",
  "name": "Human Readable Name",
  "description": "What this template does",
  "category": "backup|monitoring|cicd|database|api|devops|reporting|misc",
  "author": {
    "name": "Your Name",
    "github": "your-github-username"
  },
  "job": {
    "type": "http|shell|docker",
    "config": { ... }
  },
  "tags": ["tag1", "tag2"],
  "created_at": "YYYY-MM-DD"
}
```

3. Submit a pull request

### Template Guidelines

- Never include secrets, API keys, or credentials
- Document required environment variables
- Use lowercase kebab-case for IDs
- One job per template

## Deployment

The site is automatically deployed to GitHub Pages when changes are pushed to the `main` branch.

Manual deployment:

```bash
npm run build
# Deploy dist/ folder to your hosting provider
```

## License

MIT License - see the main [OneOff repository](https://github.com/meysam81/oneoff) for details.
