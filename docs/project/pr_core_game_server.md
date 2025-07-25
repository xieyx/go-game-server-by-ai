# Pull Request: Implement Core Game Server Modules

## Description
This PR implements the core game server modules including player, room, and game management functionality.

### Added
- Player module with player management functionality
- Room module with room management functionality
- Game module with game server functionality
- Unit tests for all core modules
- Updated CHANGELOG.md with new additions

### Implementation Details
- **Player Module**: Manages player information including ID, name, level, score, and timestamps
- **Room Module**: Manages game rooms including player capacity, status, and player management
- **Game Module**: Manages the overall game server including room and player registration, game state management

### Testing
- Added comprehensive unit tests for all core modules
- All tests pass successfully

Fixes # (issue)

## Type of change
- [x] New feature (non-breaking change which adds functionality)
- [ ] Bug fix (non-breaking change which fixes an issue)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] This change requires a documentation update

## How Has This Been Tested?
- Unit tests for player, room, and game modules
- Integration tests for game server functionality

## Checklist:
- [x] My code follows the style guidelines of this project
- [x] I have performed a self-review of my own code
- [x] I have commented my code, particularly in hard-to-understand areas
- [x] I have made corresponding changes to the documentation
- [x] My changes generate no new warnings
- [x] I have added tests that prove my fix is effective or that my feature works
- [x] New and existing unit tests pass locally with my changes
- [x] Any dependent changes have been merged and published in downstream modules
