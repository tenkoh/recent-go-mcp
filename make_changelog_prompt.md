# Go Version Changelog JSON Creation Prompt

## Objective
Create a comprehensive JSON changelog file for Go version releases that matches the established format and quality standards of existing files (go1.21.json, go1.22.json, go1.23.json).

## Input Requirements
You will be provided with:
1. **Go version number** (e.g., "1.25")
2. **Official Go release notes URL** (e.g., https://go.dev/doc/go1.25)
3. **Go language specification changes** (if applicable)
4. **Standard library documentation** for new/changed packages

## Output Format
Generate a JSON file with the following exact structure:

```json
{
  "version": "1.XX",
  "release_date": "YYYY-MM-DDTHH:MM:SSZ",
  "summary": "One-sentence compelling summary highlighting revolutionary/major features",
  "changes": [
    {
      "category": "language|runtime|toolchain|platform",
      "description": "Specific technical description with concrete details",
      "impact": "new|enhancement|performance|breaking|deprecation"
    }
  ],
  "packages": {
    "package/name": [
      {
        "function": "FunctionName (optional)",
        "description": "Detailed function/feature description",
        "impact": "new|enhancement|performance|breaking|deprecation",
        "example": "Practical code example (optional)"
      }
    ]
  }
}
```

## Content Guidelines

### Summary Section
- **One compelling sentence** (max 150 chars)
- Use power words: "revolutionary", "introduces", "major", "significant", "enhanced"
- Highlight 2-3 most impactful features
- Format: "Go 1.XX introduces [key feature], [key feature], and [performance/tooling improvement]."

### Changes Array - Information Prioritization
**INCLUDE (High Priority):**
- New language features (syntax, semantics, built-ins)
- Breaking changes (platform requirements, behavior changes)
- Significant performance improvements (with metrics when available)
- Major runtime improvements (GC, scheduler, memory management)
- Important toolchain additions (new commands, major tool enhancements)
- Platform support changes

**EXCLUDE (Low Priority):**
- Minor bug fixes without user-visible impact
- Internal refactoring without functional changes
- Documentation-only updates
- Dependency version bumps without feature impact

### Category Classification
- **language**: Syntax, semantics, built-in functions, type system
- **runtime**: GC, scheduler, memory management, performance, WASM
- **toolchain**: go command, compiler, linker, vet, fmt, build system
- **platform**: OS support, hardware architecture, system requirements

### Impact Classification
- **new**: Completely new functionality, APIs, commands
- **enhancement**: Improvements to existing functionality
- **performance**: Speed, memory, or efficiency improvements
- **breaking**: Backwards incompatible changes
- **deprecation**: Features marked for future removal

### Packages Section - Standards

**Function/Feature Selection Criteria:**
1. **Revolutionary features** - New paradigms (e.g., generics, iterators)
2. **Commonly used functionality** - High developer adoption potential
3. **Breaking changes** - Any backwards incompatible modifications
4. **Performance critical** - Significant speed/memory improvements
5. **Security relevant** - Cryptography, TLS, authentication changes

**Description Quality Standards:**
- **Specific over generic**: "Returns random int in [0, n)" vs "generates random numbers"
- **Technical precision**: Include function signatures, parameter types
- **Context provision**: Explain what problem this solves
- **Comparison context**: How it improves over previous versions

**Example Code Standards:**
- **Practical usage**: Real-world scenarios, not toy examples
- **Self-contained**: Understandable without external context
- **Concise**: 1-3 lines maximum
- **Demonstrative**: Shows the key benefit/usage pattern
- **Correct syntax**: Compilable Go code

## Information Extraction Strategy

### From Official Release Notes
1. **Scan for "Language Changes"** section - capture ALL items
2. **Runtime section** - focus on performance metrics and major improvements
3. **Standard library** - prioritize new packages and significant API additions
4. **Breaking changes** - document ALL backwards incompatible changes

### From Package Documentation
1. **New packages**: Document ALL exported functions with examples
2. **Existing packages**: Focus on new functions and breaking changes
3. **Performance improvements**: Include metrics when mentioned
4. **Deprecations**: Note replacement recommendations

### Information Validation
- **Cross-reference** package docs with release notes
- **Verify examples** are syntactically correct
- **Confirm impact levels** match actual backwards compatibility
- **Validate version numbers** and dates against official sources

## Quality Assurance Checklist

### Content Completeness
- [ ] All major language changes documented
- [ ] All new standard library packages included
- [ ] Breaking changes clearly marked
- [ ] Performance improvements quantified when possible
- [ ] New functions have usage examples

### Technical Accuracy
- [ ] Function signatures are correct
- [ ] Code examples compile and run
- [ ] Impact classifications match actual behavior
- [ ] Version requirements are accurate

### Consistency Standards
- [ ] Follows established JSON structure exactly
- [ ] Description style matches existing files
- [ ] Impact categories used consistently
- [ ] Example code style is uniform

### LLM Optimization
- [ ] Information is searchable and filterable
- [ ] Descriptions are self-contained
- [ ] Examples demonstrate practical usage
- [ ] Categories enable efficient feature discovery

## Example Quality Benchmark

**High Quality Entry:**
```json
{
  "function": "N",
  "description": "Generic random number generator - returns random value in [0, n)",
  "impact": "new",
  "example": "rand.N(100) // returns random int in [0, 100)\nrand.N(uint64(1000)) // works with any integer type"
}
```

**Low Quality Entry (Avoid):**
```json
{
  "function": "N",
  "description": "New random function",
  "impact": "new",
  "example": "rand.N(n)"
}
```

## Final Validation

Before submitting the JSON file:
1. **Parse validity**: Ensure JSON is syntactically correct
2. **Schema compliance**: Verify all required fields are present
3. **Content review**: Check against quality benchmarks
4. **Consistency check**: Compare style with existing changelog files
5. **Completeness audit**: Confirm no major features were missed

This prompt ensures reproducible, high-quality changelog files that match the established standards and provide maximum value for developers and LLM consumption.