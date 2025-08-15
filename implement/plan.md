# Implementation Plan - HTMX + Tailwind Template System

**Created**: 2025-01-15 09:00:00  
**Status**: Planning  

## Source Analysis

**Source Type**: Requirements specification for HTMX + Tailwind template architecture  
**Core Features**: Fragment-first template system with HTMX interactions and Tailwind styling  
**Target**: Multi-guild Rocket League management platform (USL as first implementation)  
**Complexity**: High - Complete template system rewrite with new architecture  

## Current State Analysis

**Existing Templates** (to be replaced):
- 18 legacy templates with inline CSS and full-page structure
- Bootstrap-based styling with duplicated CSS
- Traditional form submissions with full page reloads
- USL-specific templates in `/templates/` directory

**Current Architecture**:
- Go templates with `{{define "content"}}` pattern
- Handlers return complete HTML pages
- Basic template inheritance via `layout.html`
- No component reusability

**Target Architecture**:
- Fragment-first HTMX design
- Tailwind utility-first styling
- Guild-agnostic reusable components
- Progressive enhancement patterns

## Integration Strategy

### 1. Template Directory Structure
```
templates/
├── components/           # Reusable UI components
│   ├── form-field.html   # Input fields with validation
│   ├── user-row.html     # Individual user table row
│   ├── status-toggle.html # Active/Banned status component
│   ├── stat-card.html    # Statistics display card
│   ├── playlist-section.html # 1v1/2v2/3v3 form sections
│   ├── mmr-calculator.html # Real-time MMR display
│   ├── user-selector.html # Discord ID dropdown
│   ├── playlist-stats.html # Tracker playlist columns
│   ├── trueskill-display.html # μ/σ values display
│   └── user-trueskill-actions.html # Individual TrueSkill buttons
├── fragments/            # HTMX target fragments
│   ├── user-table.html   # Users table for swapping
│   ├── user-search.html  # Search form fragment
│   ├── user-form.html    # User creation/edit form
│   ├── user-edit-form.html # Modal edit form
│   ├── user-actions.html # Edit/Delete buttons
│   ├── delete-confirm.html # Confirmation modal
│   ├── tracker-form.html # Tracker creation form
│   ├── tracker-table.html # Trackers table
│   ├── trueskill-stats.html # TrueSkill service info
│   ├── bulk-actions.html # Bulk TrueSkill actions
│   ├── stats-grid.html   # Dashboard statistics
│   └── quick-actions.html # Navigation buttons
├── pages/                # Full page layouts
│   ├── admin-dashboard.html # Main admin page
│   ├── users.html        # Users management page
│   ├── user-form.html    # Add new user page
│   ├── trackers.html     # Trackers management page
│   ├── tracker-form.html # Add new tracker page
│   └── trueskill-dashboard.html # TrueSkill admin page
└── partials/             # Shared pieces
    ├── layout.html       # Base layout with HTMX
    ├── navigation.html   # Guild-aware navigation
    ├── modal.html        # Modal container
    └── loading.html      # Loading indicators
```

### 2. Handler Updates Required
- Modify existing handlers to return fragments instead of full pages
- Add new fragment endpoints for HTMX targets
- Implement guild context in all responses
- Add real-time calculation endpoints

### 3. Tailwind Integration
- Add Tailwind CSS to the project
- Create design system with consistent component classes
- Implement guild-aware theming system
- Remove Bootstrap and custom CSS

## Implementation Tasks

### Phase 1: Foundation Setup
- [ ] Add Tailwind CSS to the project dependencies
- [ ] Create new template directory structure
- [ ] Update base layout.html with HTMX and Tailwind
- [ ] Create guild context system for templates
- [ ] Set up component class definitions

### Phase 2: Core Components
- [ ] Build reusable form components (form-field.html, status-toggle.html)
- [ ] Create table components (user-row.html, playlist-stats.html)
- [ ] Implement status and display components (stat-card.html, trueskill-display.html)
- [ ] Build user selection and MMR calculation components

### Phase 3: User Management Templates
- [ ] Create users.html page with HTMX integration
- [ ] Build user-table.html fragment for search/filtering
- [ ] Implement user-form.html for creation/editing
- [ ] Add user action fragments (edit, delete, TrueSkill update)
- [ ] Create modal-based edit forms

### Phase 4: Tracker Management Templates
- [ ] Build trackers.html page
- [ ] Create complex tracker-form.html with playlist sections
- [ ] Implement real-time MMR calculation
- [ ] Build tracker table fragments with playlist statistics

### Phase 5: TrueSkill Management Templates
- [ ] Create trueskill-dashboard.html page
- [ ] Build TrueSkill service info fragments
- [ ] Implement individual and bulk TrueSkill action buttons
- [ ] Add progress indicators and status feedback

### Phase 6: Admin Dashboard
- [ ] Build admin-dashboard.html with statistics grid
- [ ] Create stat cards with real-time updates
- [ ] Implement quick action navigation
- [ ] Add system status monitoring

### Phase 7: Handler Integration
- [ ] Update user handlers for fragment responses
- [ ] Add HTMX endpoints for search, filtering, actions
- [ ] Implement tracker calculation endpoints
- [ ] Add TrueSkill recalculation endpoints
- [ ] Create guild context middleware

### Phase 8: Progressive Enhancement
- [ ] Add loading indicators and optimistic updates
- [ ] Implement error handling with inline messages
- [ ] Add form validation with real-time feedback
- [ ] Create accessibility improvements

### Phase 9: Guild System Integration
- [ ] Implement guild-aware routing
- [ ] Add guild theme system
- [ ] Create guild configuration templates
- [ ] Test multi-guild functionality

## API Endpoints to Implement

### Fragment Endpoints
```
GET  /{guild}/users/table           → user-table.html fragment
GET  /{guild}/users/search          → filtered user-table.html
GET  /{guild}/users/{id}/edit-form  → user-edit-form.html modal
POST /{guild}/users                 → success/error fragment
PUT  /{guild}/users/{id}            → updated user row fragment
DELETE /{guild}/users/{id}          → empty response (row removal)

GET  /{guild}/trackers/table        → tracker-table.html fragment
POST /{guild}/trackers/calculate-mmr → mmr-calculator.html fragment
POST /{guild}/trackers              → success/error fragment

POST /{guild}/users/{id}/trueskill/recalculate → trueskill-display.html
POST /{guild}/trueskill/update-all  → bulk-status fragment
GET  /{guild}/trueskill/stats       → trueskill-stats.html fragment

GET  /{guild}/admin/stats-grid      → stats-grid.html fragment
```

## Design System Classes

### Form Components
- `form-input`: "w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
- `form-select`: "w-full px-4 py-2 border border-gray-300 rounded-lg bg-white"
- `form-error`: "text-red-600 text-sm mt-1"
- `form-label`: "block text-sm font-medium text-gray-700 mb-1"

### Button Components  
- `btn-primary`: "bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 transition-colors"
- `btn-secondary`: "bg-gray-600 text-white px-4 py-2 rounded-lg hover:bg-gray-700 transition-colors"
- `btn-danger`: "bg-red-600 text-white px-4 py-2 rounded-lg hover:bg-red-700 transition-colors"
- `btn-sm`: "px-3 py-1 text-sm rounded"

### Status Components
- `status-active`: "bg-green-100 text-green-800 px-2 py-1 rounded-full text-xs font-medium"
- `status-banned`: "bg-red-100 text-red-800 px-2 py-1 rounded-full text-xs font-medium"
- `status-inactive`: "bg-gray-100 text-gray-800 px-2 py-1 rounded-full text-xs font-medium"

### Table Components
- `table-header`: "bg-gray-50 border-b border-gray-200 px-4 py-3 text-left"
- `table-row`: "border-b border-gray-100 hover:bg-gray-50 px-4 py-3"
- `table-cell`: "px-4 py-3 text-sm"

### Layout Components
- `card`: "bg-white rounded-lg shadow border p-6"
- `card-header`: "border-b border-gray-200 pb-4 mb-4"
- `stat-card`: "bg-white p-6 rounded-lg shadow border-l-4"

## Validation Checklist

- [ ] All 15 user stories implemented with HTMX interactions
- [ ] Guild-agnostic components work with theme system
- [ ] All forms work without JavaScript (progressive enhancement)
- [ ] Real-time features (search, MMR calculation, TrueSkill) functional
- [ ] Error handling with inline feedback
- [ ] Loading states and optimistic updates
- [ ] Accessibility compliance (WCAG)
- [ ] Mobile responsive design
- [ ] Performance optimization (minimal JavaScript)
- [ ] Integration with existing Go handlers
- [ ] Backward compatibility during transition
- [ ] Documentation for component usage

## Risk Mitigation

**Potential Issues**:
1. **Template compilation errors** - Test each template incrementally
2. **HTMX integration complexity** - Start with simple interactions
3. **Tailwind build process** - Integrate CSS generation pipeline
4. **Handler response format changes** - Gradual migration approach
5. **Guild context propagation** - Comprehensive middleware testing

**Rollback Strategy**:
- Keep existing templates during development
- Feature flags for new template system
- Git branches for each implementation phase
- Incremental deployment with fallback options

## Success Metrics

- [ ] Page load speed improved by >40% (eliminate full page reloads)
- [ ] CSS size reduced by >80% (eliminate duplication)
- [ ] Development velocity improved (reusable components)
- [ ] User experience enhanced (instant interactions)
- [ ] Code maintainability improved (single component system)
- [ ] Multi-guild architecture ready for expansion

## Next Steps

1. **Setup Tailwind CSS** - Add build pipeline and basic configuration
2. **Create base layout** - HTMX-enabled layout.html with guild context
3. **Build first components** - Start with simple form-field and user-row
4. **Implement user management** - Complete user CRUD with HTMX
5. **Expand to tracker and TrueSkill features** - Complex forms and calculations
6. **Deploy and test** - Validate entire system with real USL data